package botactivity

import (
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
)

func (a *ActivityBotService) GetFlow(payload types.PayloadBot) (*[]model.BotWorkflow, error) {

	botAccount, err := a.Account.GetBotByAccount(payload.MetaData.AccountId)
	if err != nil {
		return nil, err
	}
	workflow, err := a.Bot.GetBotWorkflow(botAccount.ID)
	if err != nil {
		return nil, err
	}

	return &workflow, err
}
