package util

import (
	"context"
	"fmt"
	"time"

	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

func WaitUntilCompletedOrFailAfter(ctx context.Context, apiClient *sailpoint.APIClient, taskId string, maxWaitTimeSec int64) error {
	timeoutTime := time.Now().Unix() + maxWaitTimeSec
	for {
		var lastResponse string
		status, _, _ := apiClient.Beta.TaskManagementAPI.GetTaskStatus(ctx, taskId).Execute()
		if status != nil {
			lastResponse = PrettyPrint(status)
			if status.CompletionStatus.IsSet() && status.CompletionStatus.Get() != nil && *status.CompletionStatus.Get() == "SUCCESS" {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
		if time.Now().Unix() > timeoutTime {
			return fmt.Errorf("task did not complete within %d seconds. Last response %s", maxWaitTimeSec, lastResponse)
		}
	}
}
