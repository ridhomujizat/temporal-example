package mqtt

import (
	"encoding/json"
	"fmt"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/redis"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Setup(config *Config, rds redis.IRedis) (IMqtt, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.URL)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)

	// Add reconnection settings
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Minute * 10)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Second * 30)
	opts.SetKeepAlive(time.Second * 60)

	client := &Client{
		redis: rds,
	}

	client.reconnectHandler = func(c mqtt.Client) {
		logger.Info.Println("Reconnected to MQTT broker")
		client.subscriptions.Range(func(topic, handler interface{}) bool {
			if token := c.Subscribe(topic.(string), 1, handler.(mqtt.MessageHandler)); token.Wait() && token.Error() != nil {
				logger.Error.Printf("Failed to resubscribe to topic %s: %v\n", topic, token.Error())
			} else {
				logger.Info.Printf("Resubscribed to topic: %s\n", topic)
			}
			return true
		})
	}

	opts.OnConnect = func(c mqtt.Client) {
		logger.Info.Println("Connected to MQTT broker")
		if client.reconnectHandler != nil {
			client.reconnectHandler(c)
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		logger.Error.Printf("Connection lost: %v\n", err)
	}

	opts.SetCleanSession(true)
	opts.SetWriteTimeout(time.Second * 10)
	opts.SetPingTimeout(time.Second * 60)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to broker: %v", token.Error())
	}

	client.client = mqttClient
	return client, nil
}

func (m *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) {
	m.subscriptions.Store(topic, callback)

	if token := m.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		logger.Error.Printf("Failed to subscribe to topic %s: %v\n", topic, token.Error())
		m.subscriptions.Delete(topic)
	} else {
		logger.Info.Printf("Subscribed to topic: %s\n", topic)
	}
}

func (m *Client) Publish(topic string, qos byte, retained bool, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	token := m.client.Publish(topic, qos, retained, data)
	if !token.WaitTimeout(time.Second * 10) {
		return fmt.Errorf("publish timeout for topic %s", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %v", topic, token.Error())
	}
	return nil
}

func (m *Client) Disconnect(timeout uint) {
	if m.client.IsConnected() {
		m.client.Disconnect(timeout)
	}
}

func (m *Client) AddClient(clientKey *ClientKey, clientBody *ClientBody, expr time.Duration) error {
	encBody, err := clientBody.Encrypt()
	if err != nil {
		return err
	}
	err = m.redis.Set(fmt.Sprintf(`[%q,%q,%q]`, clientKey.MountPoint, clientKey.ClientID, clientKey.Username), encBody, expr)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) ExtendTTLClient(clientKey *ClientKey) error {
	err := m.redis.Expire(fmt.Sprintf(`[%q,%q,%q]`, clientKey.MountPoint, clientKey.ClientID, clientKey.Username), 24*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) RemoveClient(clientKey *ClientKey) error {
	err := m.redis.Del(fmt.Sprintf(`[%q,%q,%q]`, clientKey.MountPoint, clientKey.ClientID, clientKey.Username))
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) Close() {
	if m.client.IsConnected() {
		m.client.Disconnect(250)
	}
}
