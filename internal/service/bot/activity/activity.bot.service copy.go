package botactivity

import (
	accountRepository "onx-outgoing-go/internal/repository/account"
	botRepository "onx-outgoing-go/internal/repository/bot"
)

type ActivityBotService struct {
	Bot     botRepository.IRepository
	Account accountRepository.IRepository
}
