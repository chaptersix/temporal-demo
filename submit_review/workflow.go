package submit_review

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Applciation struct {
		Name                       string
		Age                        int
		Email                      string
		HasCriminalHistory         bool
		CriminalHistoryDescription string
	}
	ReviewSignal struct {
		Outcome string
		Email   string
	}
)

var (
	ReviewChanName  = "Review_Signal"
	OutcomeApproved = "Approved"
	OutcomeRejected = "Rejected"
	TaskQueue       = "submit_review"
)

// SubmitAndReview workflow definition
func SubmitAndReview(ctx workflow.Context, app Applciation) (result string, err error) {
	// step 1, create new expense report
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	act := &Activities{}
	ctx1 := workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	if app.HasCriminalHistory {
		var signal ReviewSignal
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(workflow.GetSignalChannel(ctx, ReviewChanName), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &signal)
		})
		selector.Select(ctx)
		if signal.Outcome == OutcomeApproved {
			err := workflow.ExecuteActivity(ctx1, act.SendApprovalEmail, signal.Email).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to send email", "Error", err)
				return "", err
			}
		} else if signal.Outcome == OutcomeRejected {
			err := workflow.ExecuteActivity(ctx1, act.SendApprovalEmail, signal.Email).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to send email", "Error", err)
				return "", err
			}

		} else {
			return "", errors.New("unknown signal type")
		}

	} else {
		err := workflow.ExecuteActivity(ctx1, act.SendRejectedEmail, app.Email).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to send email", "Error", err)
			return "", err
		}
	}
	return "COMPLETED", nil
}
