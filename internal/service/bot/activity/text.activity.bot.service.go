package botactivity

import (
	"context"
	"fmt"
	"net/http"
	"onx-outgoing-go/internal/common/enum"
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
	"onx-outgoing-go/internal/pkg/helper"
)

func (a *ActivityBotService) Text(payload types.PayloadBot, block model.Block) (*interface{}, error) {
	accountSetting, err := a.Account.GetAccountSetting(payload.MetaData.AccountId)
	if err != nil {
		return nil, err
	}

	switch accountSetting.ChannelId {
	case enum.WHATSAPP_ID:

		return outgoingWhatsappText(a.Ctx, payload, accountSetting, block)
	default:
		return nil, nil
	}
}

func outgoingWhatsappText(ctx context.Context, payload types.PayloadBot, accountSetting model.AccountSetting, block model.Block) (*interface{}, error) {
	payloadText := types.OutgoingTextWhatsapp{
		RecipientType:    "INDIVIDUAL",
		MessagingProduct: "WHATSAPP",
		To:               payload.MetaData.UniqueId,
		Type:             "TEXT",
		Text: types.TextWhatsapp{
			Body: block.Content,
		},
	}
	baseUrl := fmt.Sprintf("%s%s", accountSetting.BaseURL, "/waba/api/messages")
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"x-key":        []string{accountSetting.Key},
		"account_id":   []string{accountSetting.Account},
	}
	resp, err := helper.HTTPRequest(&helper.HTTPRequestPayload{
		Method: enum.POST,
		URL:    baseUrl,
		Body:   payloadText,
	},
		&helper.HTTPRequestConfig{
			Headers: headers,
			Ctx:     ctx,
		})

	if err != nil {
		return &resp.Data, err
	}

	if resp.StatusCode != http.StatusOK {
		return &resp.Data, fmt.Errorf("failed to send message")
	}

	return &resp.Data, nil
}
