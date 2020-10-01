package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// RevisionClient is a client that interacts with the Gerrit "revision" REST APIs.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#revision-endpoints
type RevisionClient struct {
	*Client
}

// SetReview adds a review to a change.
func (c *RevisionClient) SetReview(ctx context.Context, changeID, revisionID string, ri *ReviewInput) error {
	var x interface{}
	return c.Call(ctx, http.MethodPost, fmt.Sprintf("/changes/%v/revisions/%v/review", changeID, revisionID), ri, &x)
}

// ReviewInput contains information for adding a review to a revision.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#review-input
type ReviewInput struct {
	Message string         `json:"message"`
	Labels  map[string]int `json:"labels"`
}
