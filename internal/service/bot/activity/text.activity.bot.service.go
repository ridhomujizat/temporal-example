package botactivity

import (
	"context"
	"fmt"
	"net/http"
	"onx-outgoing-go/internal/common/enum"
	types "onx-outgoing-go/internal/common/type"
	"onx-outgoing-go/internal/pkg/helper"
)

func (a *ActivityBotService) Text(payload types.PayloadBot) (string, error) {

	switch payload.MetaData.ChannelSources {
	case enum.WHATSAPP:
		payloadText := types.OutgoingText{
			RecipientType:    "INDIVIDUAL",
			MessagingProduct: "WHATSAPP",
			To:               payload.MetaData.UniqueId,
			Type:             "TEXT",
			Text: types.TextWhatsapp{
				Body: "Hello " + payload.Value,
			},
		}
		_, err := outGoingWhatsappText(a.Ctx, payload.MetaData.AccountId, payloadText)
		if err != nil {
			return "", err
		}
	default:
		return "Hello " + payload.Value, nil
	}

	return "Hello " + payload.Value, nil
}

func outGoingWhatsappText(ctx context.Context, account string, payload types.OutgoingText) (interface{}, error) {
	baseUrl := fmt.Sprintf("%s%s", "https://connect.infomedia.co.id", "/waba/api/messages")
	token := "8c94c57591a919a87835ad08c13c4cea2510fccf290ff98906b13f9afa7ede67f0c14699b187d9cbecefa0d071a28419"
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"x-key":        []string{token},
		"account_id":   []string{account},
	}
	resp, err := helper.HTTPRequest(&helper.HTTPRequestPayload{
		Method: enum.POST,
		URL:    baseUrl,
		Body:   payload,
	},
		&helper.HTTPRequestConfig{
			Headers: headers,
			Ctx:     ctx,
		})

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("failed to send message")
	}

	return resp.Data, nil

}
