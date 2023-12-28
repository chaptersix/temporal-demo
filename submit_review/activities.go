package submit_review

import (
	"context"
)

type Activities struct {
	// http clients and db clients etc
}

func (a *Activities) SendApprovalEmail(ctx context.Context, email string) error {
	return nil
}

func (a *Activities) SendRejectedEmail(ctx context.Context, email string) error {
	return nil
}
