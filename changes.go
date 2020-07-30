package gerrit

import (
	"context"
	"net/http"
)

// ChangeInfo contains information about a change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-info
type ChangeInfo struct {
	Project                string
	ID                     string
	UnresolvedCommentCount int              `json:"unresolved_comment_count"`
	TotalCommentCount      int              `json:"total_comment_count"`
	TrackingIDs            []TrackingIDInfo `json:"tracking_ids"`
}

// TrackingIDInfo describes a reference to an external tracking system.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#tracking-id-info
type TrackingIDInfo struct {
	System string
	ID     string
}

// ChangesClient is a client that interacts with the Gerrit "changes" REST API.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html
type ChangesClient struct {
	*Client
}

// GetChange retrieves a change.GetChange
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change
func (c *ChangesClient) GetChange(ctx context.Context, changeID string) (*ChangeInfo, error) {
	x := &ChangeInfo{}
	if err := c.Client.Call(ctx, http.MethodGet, "/changes/"+changeID, nil, x); err != nil {
		return nil, err
	}
	return x, nil
}

// ListChangeComments lists the published comments of all revisions of the change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-comments
func (c *ChangesClient) ListChangeComments(ctx context.Context, changeID string) (ChangeComments, error) {
	var x map[string][]CommentInfo
	if err := c.Client.Call(ctx, http.MethodGet, "/changes/"+changeID+"/comments", nil, &x); err != nil {
		return nil, err
	}
	return ChangeComments(x), nil
}

// ChangeComments is a mapping PATH -> CommentInfo
type ChangeComments map[string][]CommentInfo

// AccountInfo contains information about an account.
// https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#account-info
type AccountInfo struct {
	Name     string
	Email    string
	Username string
}

// CommentInfo contains information about a comment.
type CommentInfo struct {
	ID              string      `json:"id"`
	Updated         Timestamp   `json:"timestamp"`
	PatchSet        int         `json:"patch_set"`
	ChangeMessageID string      `json:"change_message_id"`
	Author          AccountInfo `json:"author"`
	InReplyTo       string      `json:"in_reply_to"`
	Message         string      `json:"message"`
	Unresolved      bool        `json:"unresolved"`
}
