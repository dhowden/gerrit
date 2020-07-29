package stream

type UnknownEventType struct {
	UnknownType string

	Data map[string]interface{}
}

// Type of the event.
func (u UnknownEventType) Type() string { return u.UnknownType }

type AssigneeChanged struct {
	Change      Change  `json:"change,omitempty"` // Change associated with the event.
	Changer     Account `json:"changer,omitempty"`
	OldAssignee Account `json:"oldAsignee,omitempty"`
}

// Type of the event.
func (AssigneeChanged) Type() string { return EventTypeAssigneeChanged }

type ChangeAbandoned struct {
	Change    Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet  PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Abandoner Account  `json:"abandoner,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

// Type of the event.
func (ChangeAbandoned) Type() string { return EventTypeChangeAbandoned }

type ChangeDeleted struct {
	Change  Change  `json:"change,omitempty"` // Change associated with the event.
	Deleter Account `json:"deleter,omitempty"`
}

// Type of the event.
func (ChangeDeleted) Type() string { return EventTypeChangeDeleted }

type ChangeMerged struct {
	Change    Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet  PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Submitter Account  `json:"submitter,omitempty"`
	NewRev    string   `json:"newRev,omitempty"`
}

// Type of the event.
func (ChangeMerged) Type() string { return EventTypeChangeMerged }

type ChangeRestored struct {
	Change   Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Restorer Account  `json:"restorer,omitempty"`
	Reason   string   `json:"reason,omitempty"`
}

// Type of the event.
func (ChangeRestored) Type() string { return EventTypeChangeRestored }

type CommentAdded struct {
	Change    Change     `json:"change,omitempty"`   // Change associated with the event.
	PatchSet  PatchSet   `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Author    Account    `json:"author,omitempty"`
	Approvals []Approval `json:"approvals,omitempty"`
	Comment   string     `json:"comment,omitempty"`
}

// Type of the event.
func (CommentAdded) Type() string { return EventTypeCommentAdded }

type DroppedOutput struct {
}

// Type of the event.
func (DroppedOutput) Type() string { return EventTypeDroppedOutput }

type HashtagsChanged struct {
	Change   Change   `json:"change,omitempty"` // Change associated with the event.
	Editor   Account  `json:"editor,omitempty"`
	Added    []string `json:"added,omitempty"`    // List of hashtags added to the change
	Removed  []string `json:"removed,omitempty"`  // List of hashtags removed from the change
	HashTags []string `json:"hashTags,omitempty"` // List of hashtags on the change after the update
}

// Type of the event.
func (HashtagsChanged) Type() string { return EventTypeHashtagsChanged }

type ProjectCreated struct {
	ProjectName string `json:"projectName,omitempty"`
	ProjectHead string `json:"projectHead,omitempty"`
}

// Type of the event.
func (ProjectCreated) Type() string { return EventTypeProjectCreated }

type PatchsetCreated struct {
	Change   Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Uploader Account  `json:"uploader,omitempty"`
}

// Type of the event.
func (PatchsetCreated) Type() string { return EventTypePatchsetCreated }

type RefUpdated struct {
	Submitter Account   `json:"submitter,omitempty"`
	RefUpdate RefUpdate `json:"refUpdate,omitempty"`
}

// Type of the event.
func (RefUpdated) Type() string { return EventTypeRefUpdated }

type ReviewerAdded struct {
	Change   Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Reviewer Account  `json:"reviewer,omitempty"`
	Adder    Account  `json:"account,omitempty"`
}

// Type of the event.
func (ReviewerAdded) Type() string { return EventTypeReviewerAdded }

type ReviewerDeleted struct {
	Change    Change     `json:"change,omitempty"`   // Change associated with the event.
	PatchSet  PatchSet   `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Reviewer  Account    `json:"reviewer,omitempty"`
	Remover   Account    `json:"remover,omitempty"`
	Approvals []Approval `json:"approvals,omitempty"`
	Comment   string     `json:"comment,omitempty"`
}

// Type of the event.
func (ReviewerDeleted) Type() string { return EventTypeReviewerDeleted }

type TopicChanged struct {
	Change   Change  `json:"change,omitempty"` // Change associated with the event.
	Changer  Account `json:"changer,omitempty"`
	OldTopic string  `json:"oldTopic,omitempty"`
}

// Type of the event.
func (TopicChanged) Type() string { return EventTypeTopicChanged }

type WIPStateChanged struct {
	Change   Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Changer  Account  `json:"changer,omitempty"`
}

// Type of the event.
func (WIPStateChanged) Type() string { return EventTypeWIPStateChanged }

type PrivateStateChanged struct {
	Change   Change   `json:"change,omitempty"`   // Change associated with the event.
	PatchSet PatchSet `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Changer  Account  `json:"changer,omitempty"`
}

// Type of the event.
func (PrivateStateChanged) Type() string { return EventTypePrivateStateChanged }

type VoteDeleted struct {
	Change    Change     `json:"change,omitempty"`   // Change associated with the event.
	PatchSet  PatchSet   `json:"patchSet,omitempty"` // PatchSet (of the Change) associated with the event.
	Reviewer  Account    `json:"reviewer,omitempty"`
	Remover   Account    `json:"remover,omitempty"`
	Approvals []Approval `json:"approvals,omitempty"`
	Comment   string     `json:"comment,omitempty"`
}

// Type of the event.
func (VoteDeleted) Type() string { return EventTypeVoteDeleted }
