package botService

import (
	"context"
	"fmt"
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func (w *Service) Init() error {
	c := w.temporal
	wf := worker.New(c, w.taskQueueName, worker.Options{})

	wf.RegisterWorkflow(Workflow)
	// wf.RegisterActivity(Activity)
	wf.RegisterActivity(&botactivity.ActivityBotService{})

	go func() {
		err := wf.Run(worker.InterruptCh())
		if err != nil {
			fmt.Println("Unable to start worker", err)
		}
	}()

	return nil
}

func (s *Service) ExecuteWorkflow(payload types.PayloadBot) (interface{}, error) {
	id := fmt.Sprintf("bot-%s-%s-%s", payload.MetaData.ChannelSources, payload.MetaData.AccountId, payload.MetaData.UniqueId)
	workflowOptions := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: s.taskQueueName,
	}

	c := s.temporal

	fmt.Println("Starting workflow")

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, Workflow, payload)
	if err != nil {
		return nil, err
	}

	fmt.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
