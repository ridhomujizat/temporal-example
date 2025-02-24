package repository

import (
	accountRepository "onx-outgoing-go/internal/repository/account"
	botRepository "onx-outgoing-go/internal/repository/bot"
)

type Repository struct {
	Account accountRepository.IRepository
	Bot     botRepository.IRepository
}
