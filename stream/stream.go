// Package stream provides tools for using Gerrit event streams.
// See https://gerrit-review.googlesource.com/Documentation/json.html for
// for futher details.
package stream

import (
	"encoding/json"
	"strconv"
	"time"
)

// UnixTime is a time.Time wrapper which decodes values
// from unix time (seconds since the Unix Epoch) used in the
// Gerrit events API.
type UnixTime time.Time

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, time.Time(ut).Unix(), 10), nil
}

func (ut *UnixTime) UnmarshalJSON(b []byte) error {
	n, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return err
	}
	*ut = UnixTime(time.Unix(n, 0))
	return nil
}

// Time returns the time.Time version of the UnixTime
// value.
func (ut *UnixTime) Time() time.Time { return time.Time(*ut) }

// Account is a Gerrit user account.
type Account struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// PatchSet refers to a specific patchset within a Change.
// https://gerrit-review.googlesource.com/Documentation/json.html#patchSet
type PatchSet struct {
	Number         int      `json:"number"`
	Revision       string   `json:"revision"`
	Parents        []string `json:"parents"`
	Ref            string   `json:"ref"`
	Uploader       Account  `json:"uploader"`
	Author         Account  `json:"author"`
	CreatedOn      UnixTime `json:"createdOn"`
	Kind           string   `json:"kind"`
	Approvals      []Approval
	Comments       []PatchsetComment
	Files          []File
	SizeInsertions int `json:"sizeInsertions"`
	SizeDeletions  int `json:"sizeDeletions"`
}

func (p *PatchSet) Accounts() []Account {
	return []Account{
		p.Uploader,
		p.Author,
	}
}

// PatchsetComment is a comment added on a patchset by a reviewer.
// https://gerrit-review.googlesource.com/Documentation/json.html#patchsetcomment
type PatchsetComment struct {
	File     string
	Line     int
	Reviewer Account
	Message  string
}

// File contains information about a patch on a file.
// https://gerrit-review.googlesource.com/Documentation/json.html#file
type File struct {
	File       string
	FileOld    string
	Type       string
	Insertions int
	Deletions  int
}

// Change represents the Gerrit change being reviewed, or that was already reviewed.
// https://gerrit-review.googlesource.com/Documentation/json.html#change
type Change struct {
	Project         string       `json:"project"`
	Branch          string       `json:"branch"`
	ID              string       `json:"id"`
	Number          int          `json:"number"`
	Subject         string       `json:"subject"`
	Owner           Account      `json:"owner"`
	URL             string       `json:"url"`
	CommitMessage   string       `json:"commitMessage"`
	HashTags        []string     `json:"hashTags"`
	CreatedOn       UnixTime     `json:"createdOn"`
	LastUpdated     UnixTime     `json:"lastUpdated"`
	Status          string       `json:"status"`
	Open            bool         `json:"open"`
	Private         bool         `json:"private"`
	WIP             bool         `json:"wip"`
	Comments        []Message    `json:"comments,omitempty"`
	TrackingIDs     []TrackingID `json:"trackingIds"`
	CurrentPatchSet PatchSet     `json:"currentPatchSet,omitEmpty"`
	PatchSets       []PatchSet   `json:"patchsets,omitempty"`
	DependsOn       Dependency   `json:"dependsOn,omitempty"`
	NeededBy        Dependency   `json:"neededBy,omitempty"`
	SubmitRecords   SubmitRecord
	AllReviewers    []Account `json:"allReviewers,omitempty"`
}

// SubmitRecord describes the submit status of a change.
// https://gerrit-review.googlesource.com/Documentation/json.html#submitRecord
type SubmitRecord struct {
	Status string
	Labels []Label
}

// Label describes a code review label for a change.
// https://gerrit-review.googlesource.com/Documentation/json.html#label
type Label struct {
	Label  string
	Status string
	By     Account
}

// Message is a comment added on a Change by a reviewer.
// https://gerrit-review.googlesource.com/Documentation/json.html#message
type Message struct {
	Timestamp UnixTime `json:"timestamp,omitempty"`
	Reviewer  Account  `json:"reviewer,omitempty"`
	Message   string   `json:"message,omitempty"`
}

// Dependency describes a change or patchset dependency.
// https://gerrit-review.googlesource.com/Documentation/json.html#dependency
type Dependency struct {
	ID                string
	Number            int
	Revision          int
	Ref               string
	IsCurrentPatchSet bool
}

// Change status values.
const (
	ChangeStatusNew       = "NEW"       // Change is still being reviewed.
	ChangeStatusMerged    = "MERGED"    // Change has been merged to its branch.
	ChangeStatusAbandoned = "ABANDONED" // Change was abandoned by its owner or administrator.
)

// TrackingID is a link to an issue tracking system.
// https://gerrit-review.googlesource.com/Documentation/json.html#trackingid
type TrackingID struct {
	System string `json:"system"`
	ID     string `json:"id"`
}

// RefUpdate contains information about a ref that was updated.
// https://gerrit-review.googlesource.com/Documentation/json.html#refUpdate
type RefUpdate struct {
	OldRev  string `json:"oldRev"`  // The old value of the ref, prior to the update.
	NewRev  string `json:"newRev"`  // The new value the ref was updated to. Zero value (0000000000000000000000000000000000000000) indicates that the ref was deleted.
	RefName string `json:"refName"` // Full ref name within project.
	Project string `json:"project"` // Project path in Gerrit.
}

// UnmarshalEvent unmarshals a JSON-encoded Gerrit event.
func UnmarshalEvent(b []byte) (*Event, error) {
	x := struct {
		Type           string
		EventCreatedOn UnixTime
	}{}

	if err := json.Unmarshal(b, &x); err != nil {
		return nil, err
	}

	var y EventType
	switch x.Type {
	case EventTypeAssigneeChanged:
		y = &AssigneeChanged{}

	case EventTypeChangeAbandoned:
		y = &ChangeAbandoned{}

	case EventTypeChangeDeleted:
		y = &ChangeDeleted{}

	case EventTypeChangeMerged:
		y = &ChangeMerged{}

	case EventTypeChangeRestored:
		y = &ChangeRestored{}

	case EventTypeCommentAdded:
		y = &CommentAdded{}

	case EventTypeDroppedOutput:
		y = &DroppedOutput{}

	case EventTypeHashtagsChanged:
		y = &HashtagsChanged{}

	case EventTypeProjectCreated:
		y = &ProjectCreated{}

	case EventTypePatchsetCreated:
		y = &PatchsetCreated{}

	case EventTypeRefUpdated:
		y = &RefUpdated{}

	case EventTypeReviewerAdded:
		y = &ReviewerAdded{}

	case EventTypeReviewerDeleted:
		y = &ReviewerDeleted{}

	case EventTypeTopicChanged:
		y = &TopicChanged{}

	case EventTypeWIPStateChanged:
		y = &WIPStateChanged{}

	case EventTypePrivateStateChanged:
		y = &PrivateStateChanged{}

	case EventTypeVoteDeleted:
		y = &VoteDeleted{}

	default:
		y = &UnknownEventType{
			UnknownType: x.Type,
		}
	}

	if err := json.Unmarshal(b, y); err != nil {
		return nil, err
	}

	return &Event{
		EventType:      y,
		EventCreatedOn: x.EventCreatedOn,
	}, nil
}

// Event represents the event reported on the stream.
// Note: attributes are used depending on the value of "type".
// https://gerrit-review.googlesource.com/Documentation/cmd-stream-events.html#events
type Event struct {
	EventType

	EventCreatedOn UnixTime `json:"eventCreatedOn"`
}

// EventType is an interface that describes specific event types
type EventType interface {
	Type() string
}

// Approval records the code review approval granted to a patch set.
// https://gerrit-review.googlesource.com/Documentation/json.html#approval
type Approval struct {
	Type        string
	Description string
	Value       string
	OldValue    string
	GrantedOn   UnixTime
	By          Account
}

// Event.Type values.
const (
	EventTypeAssigneeChanged     = "assignee-changed" // Assignee of a change has been modified.
	EventTypeChangeAbandoned     = "change-abandoned"
	EventTypeChangeDeleted       = "change-deleted"
	EventTypeChangeMerged        = "change-merged"
	EventTypeChangeRestored      = "change-restored"
	EventTypeCommentAdded        = "comment-added"
	EventTypeDroppedOutput       = "dropped-output"
	EventTypeHashtagsChanged     = "hashtags-changed"
	EventTypeProjectCreated      = "project-created"
	EventTypePatchsetCreated     = "patchset-created"
	EventTypeRefUpdated          = "ref-updated"
	EventTypeReviewerAdded       = "reviewer-added"
	EventTypeReviewerDeleted     = "reviewer-deleted"
	EventTypeTopicChanged        = "topic-changed"
	EventTypeWIPStateChanged     = "wip-state-changed"
	EventTypePrivateStateChanged = "private-state-changed"
	EventTypeVoteDeleted         = "vote-deleted"
)
