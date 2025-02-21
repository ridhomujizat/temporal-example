package botService

import (
	"encoding/json"
	"fmt"
	"onx-outgoing-go/internal/common/enum"
	types "onx-outgoing-go/internal/common/type"
	"onx-outgoing-go/internal/pkg/helper"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (s *Service) SubscribeWhatsappBot() error {
	queueName := fmt.Sprintf("%s:incoming-bot-whatsapp", helper.GetEnv("APP_TENANT"))
	fmt.Println("queueName: ", queueName)
	opts := rabbitmq.DefaultSubscribeOptions(queueName, false)
	subscriber, err := rabbitmq.NewSubscriber(s.ctx, s.rabbitmq, s.whatsappServiceBotHandler, opts)
	if err != nil {
		return err
	}

	if err := subscriber.Start(); err != nil {
		err = subscriber.Stop()
		if err != nil {
			logger.Error.Println("Failed to stop subscriber: ", err)
			return err
		}
		return err
	}

	return nil
}

func (s *Service) whatsappServiceBotHandler(msg *amqp.Delivery) (interface{}, error) {
	// logger.Debug.Println("Received message:", string(msg.Body))
	var incoming types.IncomingWhatspp

	if err := json.Unmarshal(msg.Body, &incoming); err != nil {
		logger.Error.Println("Failed to unmarshal incoming message: ", err)
		return nil, err
	}

	var payload types.PayloadBot

	payload.CustMessage = incoming
	payload.MetaData = types.MetaData{
		AccountId:       incoming.AccountId,
		UniqueId:        incoming.Contacts[0].WaId,
		ChannelPlatform: enum.SOCIOCONNECT,
		CustName:        incoming.Contacts[0].Profile.Name,
		ChannelSources:  enum.WHATSAPP,
		ChannelId:       enum.WHATSAPP_ID,
		DateTimestamp:   incoming.Messages[0].Timestamp,
	}

	switch incoming.Messages[0].Type {
	case "text":
		payload.Value = incoming.Messages[0].Text.Body
	case "location":
		payload.Value = fmt.Sprintf("%f,%f", incoming.Messages[0].Location.Latitude, incoming.Messages[0].Location.Longitude)
		payload.MetaData.IsLocation = true
		payload.MetaData.LocationPayload = types.LocationPayload{
			Latitude:  incoming.Messages[0].Location.Latitude,
			Longitude: incoming.Messages[0].Location.Longitude,
		}
	case "interactive":
		switch incoming.Messages[0].Interactive.Type {
		case "list_reply":
			payload.Value = incoming.Messages[0].Interactive.List_reply.Id
		case "button_reply":
			payload.Value = incoming.Messages[0].Interactive.Button_reply.Id
		default:
			return nil, fmt.Errorf("unsupported interactive type: %s", incoming.Messages[0].Interactive.Type)
		}
	default:
		return nil, fmt.Errorf("unsupported action: %s", incoming.Messages[0].Type)
	}

	return s.ExecuteWorkflow(payload)
}
