// Package stream provides tools for using Gerrit event streams.
// See https://gerrit-review.googlesource.com/Documentation/cmd-stream-events.html for
// for futher details.
package stream

import (
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

// Account is a representation of a Gerrit account.
type Account struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// PatchSet of a Change.
type PatchSet struct {
	Number         int      `json:"number"`
	Revision       string   `json:"revision"`
	Parents        []string `json:"parents"`
	Ref            string   `json:"ref"`
	Uploader       Account  `json:"uploader"`
	CreatedOn      UnixTime `json:"createdOn"`
	Author         Account  `json:"author"`
	Kind           string   `json:"kind"`
	SizeInsertions int      `json:"sizeInsertions"`
	SizeDeletions  int      `json:"sizeDeletions"`
}

func (p *PatchSet) Accounts() []Account {
	return []Account{
		p.Uploader,
		p.Author,
	}
}

// Change for a repo.
type Change struct {
	Project       string   `json:"project"`
	Branch        string   `json:"branch"`
	ID            string   `json:"id"`
	Number        int      `json:"number"`
	Subject       string   `json:"subject"`
	Owner         Account  `json:"owner"`
	URL           string   `json:"url"`
	CommitMessage string   `json:"commitMessage"`
	CreatedOn     UnixTime `json:"createdOn"`
	Status        string   `json:"status"`
}

type ChangeKey struct {
	ID string `json:"id"`
}

type RefUpdate struct {
	OldRev  string `json:"oldRev"`  // The old value of the ref, prior to the update.
	NewRev  string `json:"newRev"`  // The new value the ref was updated to. Zero value (0000000000000000000000000000000000000000) indicates that the ref was deleted.
	RefName string `json:"refName"` // Full ref name within project.
	Project string `json:"project"` // Project path in Gerrit.
}

// Event represents the event reported on the stream.
type Event struct {
	Uploader       Account   `json:"uploader,omitempty"`
	PatchSet       PatchSet  `json:"patchSet,omitempty"`
	Change         Change    `json:"change,omitempty"`
	Project        string    `json:"project"`
	RefName        string    `json:"refName"`
	ChangeKey      ChangeKey `json:"changeKey,omitempty"`
	Type           string    `json:"type"`
	RefUpdate      RefUpdate `json:"refUpdate,omitempty"`
	Reviewer       Account   `json:"reviewer,omitempty"`
	Remover        Account   `json:"remover,omitempty"`
	EventCreatedOn UnixTime  `json:"eventCreatedOn"`
}
