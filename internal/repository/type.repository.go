package repository

import accountRepository "onx-outgoing-go/internal/repository/account"

type Repository struct {
	Account accountRepository.IService
}
