package botactivity

import (
	"context"
	accountRepository "onx-outgoing-go/internal/repository/account"
	botRepository "onx-outgoing-go/internal/repository/bot"
)

type ActivityBotService struct {
	Ctx     context.Context
	Bot     botRepository.IRepository
	Account accountRepository.IRepository
}
