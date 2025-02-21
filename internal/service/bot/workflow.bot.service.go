package botService

import (
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"
	"time"

	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, payload types.PayloadBot) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", payload.MetaData.CustName)

	var result string
	err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Text, payload).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}
