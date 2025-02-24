package botService

import (
	"fmt"
	"onx-outgoing-go/internal/common/enum"
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

func GetNextEdgeFlow(flow model.BotWorkflowEdgesSlice, NodeID *string) *model.To {
	for _, node := range flow {
		if node.ID == *NodeID {
			return &node.To
		}
	}
	return nil
}

func GetNextEdgeFlowByChoice(flow []model.Choice, value string) *string {
	for _, node := range flow {
		if node.Value == value {
			return &node.NextEdgeID
		}
	}
	return nil
}

func Workflow(ctx workflow.Context, payload types.PayloadBot) (*types.ResultWorkflowBot, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	// logger.Info("HelloWorld workflow started", "name", payload.MetaData.CustName)

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

	resultHistory := types.ResultWorkflowBot{
		AccountId: payload.MetaData.AccountId,
		UniqueId:  payload.MetaData.UniqueId,
	}
	for currentState.CurrentFLow != "" {
		logger.Info("Current Node", "Name", currentState.FlowCurrent.Name)

		nodeState := RunningStateBlock{
			CurrentNode: currentState.FlowCurrent.Nodes.StartNodeID,
		}
		nodeState.NodeCurrent = GetEdgeFlow(currentState.FlowCurrent, nodeState.CurrentNode)

		var resultBlock []types.ResultBlockChat
		for nodeState.CurrentNode != "" {
			logger.Info("Current BLock", "Name", nodeState.NodeCurrent.Title)

			cwo := workflow.ChildWorkflowOptions{
				WorkflowID: fmt.Sprintf("%s-%s", nodeState.NodeCurrent.ID, currentState.FlowCurrent.Name),
			}
			ctx = workflow.WithChildOptions(ctx, cwo)

			var resultNode types.ResultBlockChat
			err = workflow.ExecuteChildWorkflow(ctx, WorkflowByBlock, payload, nodeState.NodeCurrent, currentState.FlowCurrent.Edges).Get(ctx, &resultNode)
			if err != nil {
				logger.Error("Failed to get workflow", "Error", err)
				return nil, err
			}

			resultBlock = append(resultBlock, resultNode)

			if resultNode.NextId.Type == "" {
				nodeState.CurrentNode = ""
				currentState.CurrentFLow = ""
				break
			}

			if resultNode.NextId.Type == "flow" {
				currentState.CurrentFLow = resultNode.NextId.WorkflowId
				break
			}

			if resultNode.NextId.Type == "node" {
				nodeState.CurrentNode = resultNode.NextId.NodeID
			}

		}

		resultHistory.ResultBlockChat = append(resultHistory.ResultBlockChat, resultBlock...)

	}

	return &resultHistory, nil
}

func WorkflowByBlock(ctx workflow.Context, payload types.PayloadBot, flow model.Node, edge model.BotWorkflowEdgesSlice) (*types.ResultBlockChat, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	// logger.Info("Hello from node", "name", flow.Title)

	result := types.ResultBlockChat{
		ID:   flow.ID,
		Node: flow.Title,
	}

	for _, block := range flow.Blocks {
		switch block.Type {
		case enum.TEXT:
			var resultBlock types.HistoryChatBot
			resultBlock.From = payload.MetaData.AccountId
			resultBlock.Type = "bot"
			resultBlock.Message = block.Content

			var resultActivity *interface{}
			err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Text, payload, block).Get(ctx, &resultActivity)
			if err != nil {
				logger.Error("Activity Fail", "Error", err)
				return nil, err
			}

			result.HistoryChat = append(result.HistoryChat, resultBlock)

			if block.NextEdgeID != nil {
				next := GetNextEdgeFlow(edge, block.NextEdgeID)
				if next != nil {
					result.NextId = *next
					break
				}
			}

			continue
		case enum.CHOICE:
			var resultBlock types.HistoryChatBot
			resultBlock.From = payload.MetaData.CustName
			resultBlock.Type = "user"
			resultBlock.Message = block.Content

			var resultActivity *interface{}
			err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Choice, payload, block).Get(ctx, &resultActivity)
			if err != nil {
				logger.Error("Activity Fail", "Error", err)
				return nil, err
			}
			result.HistoryChat = append(result.HistoryChat, resultBlock)

			var input types.PayloadBot
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(workflow.GetSignalChannel(ctx, "user_reply"),
				func(c workflow.ReceiveChannel, _ bool) {
					c.Receive(ctx, &input)
				})
			selector.Select(ctx)

			NextEdgeID := GetNextEdgeFlowByChoice(block.Choices, input.Value)
			next := GetNextEdgeFlow(edge, NextEdgeID)
			if next != nil {
				result.NextId = *next
				break
			}
			continue
		default:
			logger.Error("====Fail Nide Notfound====", "Error", "Block type not found")
			continue
		}
	}

	return &result, nil
}
