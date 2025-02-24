package botactivity

import (
	"context"
	"fmt"
	"net/http"
	"onx-outgoing-go/internal/common/enum"
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
	"onx-outgoing-go/internal/pkg/helper"
	"strings"
)

func (a *ActivityBotService) Choice(payload types.PayloadBot, block model.Block) (*interface{}, error) {
	accountSetting, err := a.Account.GetAccountSetting(payload.MetaData.AccountId)
	if err != nil {
		return nil, err
	}

	switch accountSetting.ChannelId {
	case enum.WHATSAPP_ID:
		return outgoingWhatsappChoice(a.Ctx, payload, accountSetting, block)
	default:
		return nil, nil
	}
}

func outgoingWhatsappChoice(ctx context.Context, payload types.PayloadBot, accountSetting model.AccountSetting, block model.Block) (*interface{}, error) {
	baseUrl := fmt.Sprintf("%s%s", accountSetting.BaseURL, "/waba/api/messages")
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"x-key":        []string{accountSetting.Key},
		"account_id":   []string{accountSetting.Account},
	}
	payload0utgoing := helper.HTTPRequestPayload{
		Method: enum.POST,
		URL:    baseUrl,
	}

	if block.IsDropdown != nil && *block.IsDropdown {
		payload0utgoing.Body = types.OutgoingButtonWhatsapp{
			RecipientType:    "INDIVIDUAL",
			MessagingProduct: "WHATSAPP",
			To:               payload.MetaData.UniqueId,
			Type:             "INTERACTIVE",
			Interactive: types.InteractiveWhatsapp{
				Type: "button",
				Body: types.BodyWhatsapp{Text: block.Content},
				Action: types.ActionWhatsapp{
					Buttons: mapChoicesToButtons(block.Choices),
				},
			},
		}
	} else {
		Button := "Select an option"
		payload0utgoing.Body = types.OutgoingListWhatsapp{
			RecipientType:    "INDIVIDUAL",
			MessagingProduct: "WHATSAPP",
			To:               payload.MetaData.UniqueId,
			Type:             "CHOICE",
			Interactive: types.InteractiveWhatsapp{
				Type: "list",
				Body: types.BodyWhatsapp{Text: block.Content},
				Action: types.ActionWhatsapp{
					Button: Button,
					Sections: []types.SectionWhatsapp{{
						Title: Button,
						Rows:  mapChoicesToSections(block.Choices)}},
				},
			},
		}
	}

	resp, err := helper.HTTPRequest(&payload0utgoing,
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

func mapChoicesToButtons(choices []model.Choice) []types.ButtonWhatsapp {
	var buttons []types.ButtonWhatsapp
	for _, c := range choices {
		buttons = append(buttons, types.ButtonWhatsapp{
			Type: "reply",
			Reply: types.ReplyWhatsapp{
				Title: c.Content,
				ID:    c.Value,
			},
		})
	}
	return buttons
}

func mapChoicesToSections(choices []model.Choice) []types.RowsWhatsapp {
	var sections []types.RowsWhatsapp
	for _, c := range choices {
		parts := strings.Split(c.Content, "|")
		title := strings.TrimSpace(parts[0])
		description := ""
		if len(parts) > 1 {
			description = strings.TrimSpace(parts[1])
			if len(description) > 71 {
				description = description[:71]
			}
		}
		sections = append(sections, types.RowsWhatsapp{Title: title, Description: description, ID: c.Value})
	}
	return sections
}
