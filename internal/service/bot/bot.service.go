package botService

import (
	"context"
	"fmt"
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func (w *Service) Init() error {
	c := w.temporal
	wf := worker.New(c, w.taskQueueName, worker.Options{})

	wf.RegisterWorkflow(Workflow)
	wf.RegisterWorkflow(WorkflowByBlock)
	// wf.RegisterActivity(Activity)
	wf.RegisterActivity(&botactivity.ActivityBotService{
		Bot:     w.repository.Bot,
		Account: w.repository.Account,
	})

	go func() {
		err := wf.Run(worker.InterruptCh())
		if err != nil {
			fmt.Println("Unable to start worker", err)
		}
	}()

	return nil
}

func (s *Service) ExecuteWorkflow(payload types.PayloadBot) (interface{}, error) {

	workflowID := fmt.Sprintf("bot-%s-%s-%s", payload.MetaData.ChannelSources, payload.MetaData.AccountId, payload.MetaData.UniqueId)
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: s.taskQueueName,
	}

	c := s.temporal

	err := c.SignalWorkflow(context.Background(), workflowID, "", "user_reply", payload)
	if err != nil {

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, Workflow, payload)
		if err != nil {
			return nil, err
		}
		var result *[]model.BotWorkflow
		err = we.Get(context.Background(), &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, nil
}
