package botRepository

import "onx-outgoing-go/internal/common/model"

func (r *Repository) GetBot(idBot uint) (model.BotAccount, error) {
	var botAccount model.BotAccount

	if err := r.db.GetDB().Where("id = ?", idBot).First(&botAccount).Error; err != nil {
		return botAccount, err
	}

	return botAccount, nil
}

func (r *Repository) GetBotWorkflow(idBot uint) ([]model.BotWorkflow, error) {
	var workflow []model.BotWorkflow

	if err := r.db.GetDB().Where("bot_id = ?", 1).Find(&workflow).Error; err != nil {
		return workflow, err
	}

	return workflow, nil
}

func (r *Repository) UpdateBot(bot model.BotAccount) error {
	if err := r.db.GetDB().Save(&bot).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBot(bot model.BotAccount) error {
	if err := r.db.GetDB().Create(&bot).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteBot(bot model.BotAccount) error {
	if err := r.db.GetDB().Delete(&bot).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBotWorkflowById(idBot uint, idWorkflow string) (model.BotWorkflow, error) {
	var workflow model.BotWorkflow

	if err := r.db.GetDB().Where("bot_id = ? AND id = ?", idBot, idWorkflow).First(&workflow).Error; err != nil {
		return workflow, err
	}

	return workflow, nil
}

func (r *Repository) UpdateBotWorkflow(workflow model.BotWorkflow) error {
	if err := r.db.GetDB().Save(&workflow).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBotWorkflow(workflow model.BotWorkflow) error {
	if err := r.db.GetDB().Create(&workflow).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteBotWorkflow(workflow model.BotWorkflow) error {
	if err := r.db.GetDB().Delete(&workflow).Error; err != nil {
		return err
	}

	return nil
}
