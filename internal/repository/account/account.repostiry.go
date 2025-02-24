package accountRepository

import "onx-outgoing-go/internal/common/model"

func (r *Repository) GetAccountSetting(account string) (model.AccountSetting, error) {
	var accountSetting model.AccountSetting

	if err := r.db.GetDB().Where("account = ?", account).First(&accountSetting).Error; err != nil {
		return accountSetting, err
	}

	return accountSetting, nil
}

func (r *Repository) UpdateAccountSetting(account string, setting model.AccountSetting) error {
	if err := r.db.GetDB().Model(&model.AccountSetting{}).Where("account = ?", account).Updates(setting).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateAccountSetting(setting model.AccountSetting) error {
	if err := r.db.GetDB().Create(&setting).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAccountSetting(account string) error {
	if err := r.db.GetDB().Where("account = ?", account).Delete(&model.AccountSetting{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBotByAccount(account string) (*model.BotAccount, error) {
	var accountSetting model.AccountSetting

	if err := r.db.GetDB().Where("account = ?", account).First(&accountSetting).Error; err != nil {
		return nil, err
	}

	var botAccount model.BotAccount

	if err := r.db.GetDB().Where("id = ?", accountSetting.BotID).First(&botAccount).Error; err != nil {
		return nil, err
	}

	return &botAccount, nil
}
