package accountRepository

import (
	"context"
	"encoding/json"
	"fmt"
	accountTypes "onx-outgoing-go/internal/common"
	"time"

	"onx-outgoing-go/internal/pkg/helper"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
)

type Service struct {
	ctx      context.Context
	redis    redis.IRedis
	rabbit   *rabbitmq.ConnectionManager
	publiser *rabbitmq.Publisher
}

type DataGetAccount struct {
	ChannelId int    `json:"channel_id"`
	AccountId string `json:"account_id"`
}

type IService interface {
	GetAccount(ChannelId int, AccountId string) (*accountTypes.SettingAccount, error)
}

func NewService(ctx context.Context, redis redis.IRedis, rabbit *rabbitmq.ConnectionManager, publisher *rabbitmq.Publisher) IService {
	return &Service{
		ctx:      ctx,
		redis:    redis,
		rabbit:   rabbit,
		publiser: publisher,
	}
}

func (s *Service) GetAccount(ChannelId int, AccountId string) (*accountTypes.SettingAccount, error) {
	key := fmt.Sprintf("%s:detailChannel:%d:%s", helper.GetEnv("APP_TENANT"), ChannelId, AccountId)
	account, err := s.getAccountFromRedis(key)

	if err != nil {
		return nil, err
	}

	if account != nil {
		return account, nil
	}

	account, err = s.getAccountFromRpcOmnix(ChannelId, AccountId)
	fmt.Println("Account", account, err)

	if err != nil {
		return nil, err
	}

	if account != nil {
		fmt.Println("Account", account)
		if err = s.redis.Set(key, account, 24*time.Hour); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func (s *Service) getAccountFromRedis(key string) (*accountTypes.SettingAccount, error) {
	accountStr, err := s.redis.Get(key)
	if err != nil {
		return nil, err
	}

	if accountStr != "" {
		var account accountTypes.SettingAccount
		if err := json.Unmarshal([]byte(accountStr), &account); err != nil {
			return nil, err
		}
		return &account, nil
	}

	return nil, nil
}

func (s *Service) getAccountFromRpcOmnix(ChannelId int, AccountId string) (*accountTypes.SettingAccount, error) {
	response, err := s.rabbit.GetRPCMicroserviceOmnix("account.byAccount", DataGetAccount{
		ChannelId: ChannelId,
		AccountId: AccountId,
	})

	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, err
	}

	var account accountTypes.SettingAccount
	if err := json.Unmarshal(jsonResponse, &account); err != nil {
		return nil, err
	}
	return &account, nil

}
