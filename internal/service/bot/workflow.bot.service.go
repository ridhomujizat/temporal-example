package botService

import (
	"fmt"
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"
	"time"

	"go.temporal.io/sdk/workflow"
)

type RunningStateFlow struct {
	CurrentFLow string
	FlowCurrent model.BotWorkflow
}

type RunningStateBlock struct {
	CurrentNode string
	NodeCurrent model.Node
}

func GetNodeFlow(workFlow []model.BotWorkflow, flowName string) model.BotWorkflow {
	for _, flow := range workFlow {
		if flow.Name == flowName {
			return flow
		}
	}
	return model.BotWorkflow{}
}

func GetEdgeFlow(flow model.BotWorkflow, NodeID string) model.Node {
	for _, node := range flow.Nodes.Nodes {
		if node.ID == NodeID {
			return node
		}
	}
	return model.Node{}
}

func Workflow(ctx workflow.Context, payload types.PayloadBot) (*[]model.BotWorkflow, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", payload.MetaData.CustName)

	//Get Workfloww
	var workFlow []model.BotWorkflow
	err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).GetFlow, payload).Get(ctx, &workFlow)
	if err != nil {
		logger.Error("Failed to get workflow", "Error", err)
		return nil, err
	}

	currentState := RunningStateFlow{
		CurrentFLow: "main",
	}
	currentState.FlowCurrent = GetNodeFlow(workFlow, currentState.CurrentFLow)

	for currentState.FlowCurrent.Name != "" {
		logger.Info("Current Node", "Name", currentState.FlowCurrent.Name)

		nodeState := RunningStateBlock{
			CurrentNode: currentState.FlowCurrent.Nodes.StartNodeID,
		}

		nodeState.NodeCurrent = GetEdgeFlow(currentState.FlowCurrent, nodeState.CurrentNode)

		for nodeState.CurrentNode != "" {
			logger.Info("Current Node", "Name", nodeState.NodeCurrent.Title)

			cwo := workflow.ChildWorkflowOptions{
				WorkflowID: fmt.Sprintf("%s-%s", nodeState.NodeCurrent.ID, currentState.FlowCurrent.Name),
			}
			ctx = workflow.WithChildOptions(ctx, cwo)

			var result interface{}
			err = workflow.ExecuteChildWorkflow(ctx, WorkflowByBlock, payload, nodeState.NodeCurrent).Get(ctx, &result)
			if err != nil {
				logger.Error("Failed to get workflow", "Error", err)
				return nil, err
			}
		}

	}

	return &workFlow, nil
}

func WorkflowByBlock(ctx workflow.Context, payload types.PayloadBot, flow model.Node) (*interface{}, error) {

	return nil, nil
}
