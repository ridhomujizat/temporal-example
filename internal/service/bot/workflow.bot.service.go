package botService

import (
	"fmt"
	"onx-outgoing-go/internal/common/enum"
	"onx-outgoing-go/internal/common/model"
	types "onx-outgoing-go/internal/common/type"
	botactivity "onx-outgoing-go/internal/service/bot/activity"
	"time"

	"go.temporal.io/sdk/temporal"
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

func GetNextEdgeFlow(flow model.BotWorkflowEdgesSlice, NodeID string) *model.To {
	for _, node := range flow {
		if node.ID == NodeID {
			return &node.To
		}
	}
	return nil
}

func GetNextEdgeFlowByChoice(flow []model.Choice, value string) string {
	for _, node := range flow {
		if node.Value == value {
			return node.NextEdgeID
		}
	}
	return ""
}

func Workflow(ctx workflow.Context, payload types.PayloadBot) (*types.ResultWorkflowBot, error) {
	// Set retry policy for activities
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Bot workflow started", "account", payload.MetaData.AccountId)

	// Get Workflow with error handling
	var workFlow []model.BotWorkflow
	err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).GetFlow, payload).Get(ctx, &workFlow)
	if err != nil {
		logger.Error("Failed to get workflow", "Error", err)
		return &types.ResultWorkflowBot{
			AccountId: payload.MetaData.AccountId,
			UniqueId:  payload.MetaData.UniqueId,
			Error:     fmt.Sprintf("Failed to fetch workflow: %v", err),
		}, err
	}

	if len(workFlow) == 0 {
		return &types.ResultWorkflowBot{
			AccountId: payload.MetaData.AccountId,
			UniqueId:  payload.MetaData.UniqueId,
			Error:     "No workflow found",
		}, fmt.Errorf("no workflow found")
	}

	currentState := RunningStateFlow{
		CurrentFLow: "main",
	}
	currentState.FlowCurrent = GetNodeFlow(workFlow, currentState.CurrentFLow)

	// Check if main flow exists
	if currentState.FlowCurrent.Name == "" {
		return &types.ResultWorkflowBot{
			AccountId: payload.MetaData.AccountId,
			UniqueId:  payload.MetaData.UniqueId,
			Error:     "Main flow not found",
		}, fmt.Errorf("main flow not found")
	}

	resultHistory := types.ResultWorkflowBot{
		AccountId: payload.MetaData.AccountId,
		UniqueId:  payload.MetaData.UniqueId,
	}

	// Main flow loop
	for currentState.CurrentFLow != "" {
		logger.Info("Processing flow", "name", currentState.FlowCurrent.Name)

		nodeState := RunningStateBlock{
			CurrentNode: currentState.FlowCurrent.Nodes.StartNodeID,
		}

		var resultBlock []types.ResultBlockChat
		// Node processing loop
		for nodeState.CurrentNode != "" {
			nodeState.NodeCurrent = GetEdgeFlow(currentState.FlowCurrent, nodeState.CurrentNode)
			if nodeState.NodeCurrent.ID == "" {
				logger.Error("Node not found", "nodeID", nodeState.CurrentNode)
				break
			}

			logger.Info("Processing node", "title", nodeState.NodeCurrent.Title)

			cwo := workflow.ChildWorkflowOptions{
				WorkflowID: fmt.Sprintf("%s-%s", nodeState.NodeCurrent.ID, currentState.FlowCurrent.Name),
			}
			ctx = workflow.WithChildOptions(ctx, cwo)

			flow := nodeState.NodeCurrent
			edge := currentState.FlowCurrent.Edges
			resultNode := types.ResultBlockChat{
				ID:   flow.ID,
				Node: flow.Title,
			}

			for _, block := range flow.Blocks {
				switch block.Type {
				case enum.TEXT:
					var resultActivity *interface{}
					err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Text, payload, block).Get(ctx, &resultActivity)
					if err != nil {
						logger.Error("Text activity failed", "error", err)
						resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
							From:    payload.MetaData.AccountId,
							Type:    "bot",
							Message: "Sorry, I'm having trouble processing your request.",
						})
						continue
					}

					resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
						From:    payload.MetaData.AccountId,
						Type:    "bot",
						Message: block.Content,
					})

					if block.NextEdgeID != nil {
						next := GetNextEdgeFlow(edge, *block.NextEdgeID)
						if next != nil {
							resultNode.NextId = *next
							break
						} else {
							currentState.CurrentFLow = ""
						}
					}

					continue

				case enum.CHOICE:
					var resultActivity *interface{}
					err := workflow.ExecuteActivity(ctx, (*botactivity.ActivityBotService).Choice, payload, block).Get(ctx, &resultActivity)
					if err != nil {
						logger.Error("Choice activity failed", "error", err)
						resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
							From:    payload.MetaData.AccountId,
							Type:    "bot",
							Message: "Sorry, I'm having trouble processing your request.",
						})
						continue
					}

					resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
						From:    payload.MetaData.AccountId,
						Type:    "bot",
						Message: block.Content,
					})

					// Use selector with timeout for user input
					var input types.PayloadBot
					userReplyReceived := false

					// Wait for user reply with a timeout
					selector := workflow.NewSelector(ctx)
					selector.AddReceive(workflow.GetSignalChannel(ctx, "user_reply"),
						func(c workflow.ReceiveChannel, _ bool) {
							c.Receive(ctx, &input)
							userReplyReceived = true
						})

					// Add timeout (e.g., 2 minutes)
					timerFuture := workflow.NewTimer(ctx, time.Minute*2)
					selector.AddFuture(timerFuture, func(f workflow.Future) {
						// Timer fired, no signal received
						userReplyReceived = false
					})

					selector.Select(ctx)

					// Handle timeout case
					if !userReplyReceived {
						logger.Info("User reply timeout")
						resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
							From:    payload.MetaData.AccountId,
							Type:    "bot",
							Message: "I haven't heard from you in a while. This conversation will end now.",
						})

						// End the conversation
						nodeState.CurrentNode = ""
						currentState.CurrentFLow = ""
						break
					}

					// Process user input
					resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
						From:    input.MetaData.CustName,
						Type:    "user",
						Message: input.Value,
					})

					NextEdgeID := GetNextEdgeFlowByChoice(block.Choices, input.Value)

					// Handle invalid choice
					if NextEdgeID == "" {
						logger.Info("Invalid choice selected", "value", input.Value)
						resultNode.HistoryChat = append(resultNode.HistoryChat, types.HistoryChatBot{
							From:    payload.MetaData.AccountId,
							Type:    "bot",
							Message: "Sorry, that's not a valid option. Please try again.",
						})

						// Repeat the same node (don't move forward)
						continue
					}

					logger.Info("Choice selected", "value", input.Value, "nextEdge", NextEdgeID)

					next := GetNextEdgeFlow(edge, NextEdgeID)
					if next != nil {
						resultNode.NextId = *next
						break
					}
					continue

				default:
					logger.Error("Unknown block type", "type", block.Type)
					continue
				}
			}

			logger.Info("Navigation", "nextType", resultNode.NextId.Type, "nextNode", resultNode.NextId.NodeID)
			resultBlock = append(resultBlock, resultNode)

			if resultNode.NextId.Type == "" {
				nodeState.CurrentNode = ""
				currentState.CurrentFLow = ""
				break
			}

			if resultNode.NextId.Type == "flow" {
				currentState.CurrentFLow = resultNode.NextId.WorkflowId
				// Update current flow with the new flow data
				currentState.FlowCurrent = GetNodeFlow(workFlow, currentState.CurrentFLow)
				if currentState.FlowCurrent.Name == "" {
					logger.Error("Target flow not found", "flowId", currentState.CurrentFLow)
					currentState.CurrentFLow = ""
				}
				break
			}

			if resultNode.NextId.Type == "node" {
				nodeState.CurrentNode = resultNode.NextId.NodeID
			}
		}

		resultHistory.ResultBlockChat = append(resultHistory.ResultBlockChat, resultBlock...)
	}

	logger.Info("Workflow completed", "account", payload.MetaData.AccountId)
	return &resultHistory, nil
}
