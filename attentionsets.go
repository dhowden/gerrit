package gerrit

import (
	"context"
	"net/http"
)

// The AttentionSetInfo entity contains details of users that are in the attention set.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#attention-set-info
type AttentionSetInfo struct {
	Account    AccountInfo `json:"account"`     // AccountInfo entity.
	LastUpdate Timestamp   `json:"last_update"` // The timestamp of the last update.
	Reason     string      `json:"reason"`      // The reason of for adding or removing the user.
}

type AttentionSetClient struct {
	*Client
}

// GetAttentionSet fetches all users that are currently in the attention set.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-attention-set
func (c *AttentionSetClient) GetAttentionSet(ctx context.Context, changeID string) ([]AttentionSetInfo, error) {
	x := []AttentionSetInfo{}
	if err := c.Client.Call(ctx, http.MethodGet, "/changes/"+changeID+"/attention", nil, &x); err != nil {
		return nil, err
	}
	return x, nil
}
