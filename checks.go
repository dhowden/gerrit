package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// For more details on the checks JSON API:
// https://gerrit.googlesource.com/plugins/checks/+/refs/heads/stable-3.2/resources/Documentation/

// CheckerCreateInput contains information for creating a checker.
type CheckerCreateInput struct {
	UUID        string `json:"uuid"`                  // The UUID of the checker.
	Name        string `json:"name"`                  // The name of the checker.
	Description string `json:"description,omitempty"` // The description of the checker.
	URL         string `json:"url,omitempty"`         // The URL of the checker.
	Repository  string `json:"repository"`            // The (exact) name of the repository for which the checker applies.
	Status      string `json:"status,omitempty"`      // The status of the checker; one of ENABLED or DISABLED.
	Blocking    string `json:"blocking,omitempty"`    // A list of conditions that describe when the checker should block change submission.
	Query       string `json:"query,omitempty"`       // A query that limits changes for which the checker is relevant.
}

// CheckerInfo describes a checker.
type CheckerInfo struct {
	UUID        string    `json:"uuid"`                  // The UUID of the checker.
	Name        string    `json:"name"`                  // The name of the checker.
	Description string    `json:"description,omitempty"` // The description of the checker.
	URL         string    `json:"url,omitempty"`         // The URL of the checker.
	Repository  string    `json:"repository"`            // The (exact) name of the repository for which the checker applies.
	Status      string    `json:"status"`                // The status of the checker; one of ENABLED or DISABLED.
	Blocking    []string  `json:"blocking"`              // A list of conditions that describe when the checker should block change submission.
	Query       string    `json:"query,omitempty"`       // A query that limits changes for which the checker is relevant.
	Created     time.Time `json:"created"`               // The timestamp of when the checker was created.
	Updated     time.Time `json:"updated"`               // The timestamp of when the checker was last updated.
}

// CheckInfo describes a check.
type CheckInfo struct {
	Repository         string     `json:"repository"`                    // The repository name that this check applies to.
	ChangeNumber       int        `json:"change_number"`                 // The change number that this check applies to.
	PatchSetID         int        `json:"patch_set_id"`                  // The patch set that this check applies to.
	CheckerUUID        string     `json:"checker_uuid"`                  // The UUID of the checker that reported this check.
	State              CheckState `json:"state"`                         // The state as string-serialized form of CheckState
	Message            string     `json:"message,omitempty"`             //	Short message explaining the check state.
	URL                string     `json:"url,omitempty"`                 //	A fully-qualified URL pointing to the result of the check on the checker’s infrastructure.
	Started            Timestamp  `json:"started,omitempty"`             //	The timestamp of when the check started processing.
	Finished           Timestamp  `json:"finished,omitempty"`            //	The timestamp of when the check finished processing.
	Created            Timestamp  `json:"created"`                       // The timestamp of when the check was created.
	Updated            Timestamp  `json:"updated"`                       // The timestamp of when the check was last updated.
	CheckerName        string     `json:"checker_name,omitempty"`        //	The name of the checker that produced this check.  Only set if checker details are requested.
	CheckerStatus      string     `json:"checker_status,omitempty"`      //	The status of the checker that produced this check.  Only set if checker details are requested.
	Blocking           []string   `json:"blocking,omitempty"`            //	Set of blocking conditions that apply to this checker.  Only set if checker details are requested.
	CheckerDescription string     `json:"checker_description,omitempty"` //	The description of the checker that reported this check.
}

// CheckInput contains information for creating or updating a check.
type CheckInput struct {
	CheckerUUID   string     `json:"checker_uuid,omitempty"`   //	The UUID of the checker. Must be specified for check creation. Optional only if updating a check and referencing the checker using the UUID in the URL.
	State         CheckState `json:"state,omitempty"`          //	The state as string-serialized form of CheckState
	Message       string     `json:"message,omitempty"`        //	Short message explaining the check state.
	URL           string     `json:"url,omitempty"`            //	A fully-qualified URL pointing to the result of the check on the checker’s infrastructure.
	Started       *Timestamp `json:"started,omitempty"`        //	The timestamp of when the check started processing.
	Finished      *Timestamp `json:"finished,omitempty"`       //	The timestamp of when the check finished processing.
	Notify        string     `json:"notify,omitempty"`         //	Notify handling that defines to whom email notifications should be sent when the combined check state changes due to posting this check. Allowed values are NONE, OWNER, OWNER_REVIEWERS and ALL. If not set, the default is ALL if the combined check state is updated to either SUCCESSFUL or NOT_RELEVANT, otherwise the default is OWNER. Regardless of this setting there are no email notifications for posting checks on non-current patch sets.
	NotifyDetails string     `json:"notify_details,omitempty"` //	Additional information about whom to notify when the combined check state changes due to posting this check as a map of recipient type to NotifyInfo entity. Regardless of this setting there are no email notifications for posting checks on non-current patch sets.
}

// ChecksClient is a client for interating with the Gerrit Checks API.
type ChecksClient struct {
	*Client
}

// Timestamp is a time.Time wrapper which decodes values
// in the layout yyyy-mm-dd hh:mm:ss.fffffffff (in UTC).
//
// See https://gerrit-review.googlesource.com/Documentation/rest-api.html#timestamp
type Timestamp time.Time

const timestampLayout = "2006-01-02 15:04:05.999999999"

func (ts Timestamp) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(timestampLayout))
	b = time.Time(ts).UTC().AppendFormat(b, timestampLayout)
	return b, nil
}

func (ts *Timestamp) UnmarshalText(b []byte) error {
	if len(b) != len(timestampLayout) {
		return fmt.Errorf("unknown date format %q", b)
	}
	t, err := time.Parse(timestampLayout, string(b))
	if err != nil {
		return err
	}
	*ts = Timestamp(t)
	return nil
}

// Time returns the time.Time version of the Timestamp
// value.
func (ts *Timestamp) Time() time.Time { return time.Time(*ts) }

// CheckablePatchSetInfo describes a patch set for which checks are pending.
type CheckablePatchSetInfo struct {
	Repository   string `json:"repository"`    // The repository name that this pending check applies to.
	ChangeNumber int    `json:"change_number"` // The change number that this pending check applies to.
	PatchSetID   int    `json:"patch_set_id"`  // The ID of the patch set that this pending check applies to.
}

// PendingCheckInfo describes a pending check.
type PendingCheckInfo struct {
	State CheckState `json:"state"` // State of the check.
}

// PendingChecksInfo describes the pending checks on patch set.
type PendingChecksInfo struct {
	PatchSet      CheckablePatchSetInfo       `json:"patch_set"`      // The patch set for checks are pending as CheckablePatchSetInfo entity.
	PendingChecks map[string]PendingCheckInfo `json:"pending_checks"` // The checks that are pending for the patch set as checker UUID to PendingCheckInfo entity.
}

// CheckState represents the state of a check.
type CheckState string

// CheckState values.
const (
	StateNotStarted  CheckState = "NOT_STARTED"
	StateFailed      CheckState = "FAILED"
	StateScheduled   CheckState = "SCHEDULED"
	StateRunning     CheckState = "RUNNING"
	StateSuccessful  CheckState = "SUCCESSFUL"
	StateNotRelevant CheckState = "NOT_RELEVANT"
)

var validCheckStates = []CheckState{
	StateNotStarted,
	StateFailed,
	StateScheduled,
	StateRunning,
	StateSuccessful,
	StateNotRelevant,
}

func (c *CheckState) UnmarshalText(b []byte) error {
	s := CheckState(b)
	for _, x := range validCheckStates {
		if x == s {
			*c = s
			return nil
		}
	}
	return fmt.Errorf("invalid check state: %q", b)
}

const (
	pendingQuery    = "query=scheme:test+(state:NOT_STARTED+OR+state:SCHEDULED)"
	notStartedQuery = "query=scheme:test+state:NOT_STARTED"
)

func (c *ChecksClient) Pending(ctx context.Context) ([]PendingChecksInfo, error) {
	var resp []PendingChecksInfo
	if err := c.Client.Call(ctx, http.MethodGet, "/plugins/checks/checks.pending/?"+pendingQuery, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ChecksClient) NotStarted(ctx context.Context) ([]PendingChecksInfo, error) {
	var resp []PendingChecksInfo
	if err := c.Client.Call(ctx, http.MethodGet, "/plugins/checks/checks.pending/?"+notStartedQuery, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ChecksClient) checkURL(changeNumber, patchSetID int) string {
	return fmt.Sprintf("/changes/%d/revisions/%d/checks", changeNumber, patchSetID)
}

func (c *ChecksClient) List(ctx context.Context, changeNumber, patchSetID int) ([]CheckInfo, error) {
	var resp []CheckInfo
	if err := c.Client.Call(ctx, http.MethodGet, c.checkURL(changeNumber, patchSetID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ChecksClient) updateCheck(ctx context.Context, changeNumber, patchSetID int, req *CheckInput) (CheckInfo, error) {
	var resp CheckInfo
	if err := c.Client.Call(ctx, http.MethodPost, c.checkURL(changeNumber, patchSetID), req, &resp); err != nil {
		return CheckInfo{}, err
	}
	return resp, nil
}

func (c *ChecksClient) Start(ctx context.Context, uuid string, changeNumber, patchSetID int, state CheckState, logURL string) (CheckInfo, error) {
	started := Timestamp(time.Now())
	req := &CheckInput{
		CheckerUUID: uuid,
		State:       state,
		Started:     &started,
		URL:         logURL,
	}
	return c.updateCheck(ctx, changeNumber, patchSetID, req)
}

func (c *ChecksClient) Update(ctx context.Context, uuid string, changeNumber, patchSetID int, state CheckState, logURL string) (CheckInfo, error) {
	req := &CheckInput{
		CheckerUUID: uuid,
		State:       state,
		URL:         logURL,
	}
	return c.updateCheck(ctx, changeNumber, patchSetID, req)
}

func (c *ChecksClient) Finish(ctx context.Context, uuid string, changeNumber, patchSetID int, state CheckState) (CheckInfo, error) {
	finished := Timestamp(time.Now())
	req := &CheckInput{
		CheckerUUID: uuid,
		State:       state,
		Finished:    &finished,
	}
	return c.updateCheck(ctx, changeNumber, patchSetID, req)
}
