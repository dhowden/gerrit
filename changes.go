package gerrit

import (
	"context"
	"net/http"
	"net/url"
)

// ChangeInfo contains information about a change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-info
type ChangeInfo struct {
	Project                string                      `json:"project"`
	ID                     string                      `json:"id"`
	ChangeID               string                      `json:"change_id"`
	UnresolvedCommentCount int                         `json:"unresolved_comment_count"`
	TotalCommentCount      int                         `json:"total_comment_count"`
	TrackingIDs            []TrackingIDInfo            `json:"tracking_ids"`
	Messages               []ChangeMessageInfo         `json:"messages"`
	Subject                string                      `json:"subject"`
	Branch                 string                      `json:"branch"`
	Created                Timestamp                   `json:"created"`
	Updated                Timestamp                   `json:"updated"`
	Submitted              Timestamp                   `json:"submitted"`
	Owner                  AccountInfo                 `json:"owner"`
	Number                 int                         `json:"_number"`
	Reviewers              map[string][]AccountInfo    `json:"reviewers"`
	Revisions              map[string]RevisionInfo     `json:"revisions"`
	AttentionSet           map[string]AttentionSetInfo `json:"attention_set"`
	Submittable            bool                        `json:"submittable"` // Only set if requested via SUBMITTABLE option.
}

// RevisionInfo contains information about a revision.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#revision-info
type RevisionInfo struct {
	Number   int `json:"_number"`
	Commit   CommitInfo
	Created  Timestamp
	Uploader AccountInfo
}

// CommitInfo contains information about a commit.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#commit-info
type CommitInfo struct {
	Parents []CommitInfo
	Subject string
	Message string
}

// ChangeMessageInfo contains information about a message attached to a change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-message-info
type ChangeMessageInfo struct {
	ID             string       // The ID of the message.
	Author         *AccountInfo // Author of the message as an AccountInfo entity.
	RealAuthor     *AccountInfo
	Date           Timestamp
	Message        string
	RevisionNumber int `json:"_revision_number,omitempty"` // Which patchset (if any) generated this message.
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

// GetChange retrieves a change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change
func (c *ChangesClient) GetChange(ctx context.Context, changeID string, opts ...string) (*ChangeInfo, error) {
	query := ""
	if len(opts) > 0 {
		v := url.Values{"o": opts}
		query = "?" + v.Encode()
	}

	x := &ChangeInfo{}
	if err := c.Client.Call(ctx, http.MethodGet, "/changes/"+changeID+query, nil, x); err != nil {
		return nil, err
	}
	return x, nil
}

// ListChangeComments lists the published comments of all revisions of the change.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-comments
func (c *ChangesClient) ListChangeComments(ctx context.Context, changeID string, opts ...string) (ChangeComments, error) {
	query := ""
	if len(opts) > 0 {
		v := url.Values{"o": opts}
		query = "?" + v.Encode()
	}

	var x map[string][]CommentInfo
	if err := c.Client.Call(ctx, http.MethodGet, "/changes/"+changeID+"/comments"+query, nil, &x); err != nil {
		return nil, err
	}
	return ChangeComments(x), nil
}

// ChangeComments is a mapping PATH -> CommentInfo.
type ChangeComments map[string][]CommentInfo

// AccountInfo contains information about an account.
// https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#account-info
type AccountInfo struct {
	Name     string
	Email    string
	Username string
}

// CommentInfo contains information about a comment.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#comment-info
type CommentInfo struct {
	ID              string       `json:"id"`
	Updated         Timestamp    `json:"updated"`
	PatchSet        int          `json:"patch_set"`
	Path            string       `json:"path"`
	Line            int          `json:"line"`
	Range           CommentRange `json:"range"`
	ChangeMessageID string       `json:"change_message_id"`
	Author          AccountInfo  `json:"author"`
	InReplyTo       string       `json:"in_reply_to"`
	Message         string       `json:"message"`
	Unresolved      bool         `json:"unresolved"`
}

// CommentRange describes the range of an inline comment.
//
// The comment range is a range from the start position, specified by
// start_line and start_character, to the end position, specified by
// end_line and end_character. The start position is inclusive and the
// end position is exclusive.
//
// So, a range over part of a line will have start_line equal to end_line;
// however a range with end_line set to 5 and end_character equal to 0 will
// not include any characters on line 5,
type CommentRange struct {
	StartLine      int // Start line number of the range (1-based).
	StartCharacter int // Character position in the start line (0-based).
	EndLine        int // End line number of the range (1-based).
	EndCharacter   int // Character position in the end line (0-based).
}
