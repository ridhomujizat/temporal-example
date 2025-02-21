package botService

import (
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"
	"time"

	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, payload types.PayloadBot) (*types.ResultWorkflowBot, error) {
	// workflow setup
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("Bot", "name", payload.MetaData.CustName)

	// data processing
	var resultWorkflowBot types.ResultWorkflowBot

	// initial data
	resultWorkflowBot.AccountId = payload.MetaData.AccountId
	resultWorkflowBot.UniqueId = payload.MetaData.UniqueId
	resultWorkflowBot.HostoryChat = append(resultWorkflowBot.HostoryChat, types.HistoryChatBot{
		From:          payload.MetaData.CustName,
		Type:          "text",
		Message:       payload.Value,
		DateTimestamp: payload.MetaData.DateTimestamp,
	})

	// activity Text
	var result types.HistoryChatBot
	err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Text, payload).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return &resultWorkflowBot, err
	}
	resultWorkflowBot.HostoryChat = append(resultWorkflowBot.HostoryChat, result)

	var input types.PayloadBot
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(workflow.GetSignalChannel(ctx, "user_reply"),
		func(c workflow.ReceiveChannel, _ bool) {
			c.Receive(ctx, &input)
		})
	selector.Select(ctx)

	// err = workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).End, payload, input.Value).Get(ctx, &result)
	// if err != nil {
	// 	logger.Error("Activity failed.", "Error", err)
	// 	return &resultWorkflowBot, err
	// }

	logger.Info("HelloWorld workflow completed.", "result", result)

	return &resultWorkflowBot, nil
}
