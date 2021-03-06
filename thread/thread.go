package thread

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/dhowden/gerrit"
)

// Summary of a change.
type Summary struct {
	ChangeID string
	Project  string
	Branch   string

	Subject string

	LatestCommitMessage string
	AllReviewers        []gerrit.AccountInfo
	ActiveReviewers     []gerrit.AccountInfo
	CCed                []gerrit.AccountInfo

	Created   time.Time
	Updated   time.Time
	Submitted time.Time

	Comments           int
	UnresolvedComments int

	Threads []Thread
}

// Thread of comments.
type Thread struct {
	s *Summary

	Path     string
	Line     int
	PatchSet int
	Authors  []gerrit.AccountInfo
	Message  string

	LastComment gerrit.CommentInfo
}

func (t *Thread) URL() string {
	return fmt.Sprintf("/c/%s/+/%s/%d/%v#%d", t.s.Project, t.s.ChangeID, t.PatchSet, t.Path, t.Line)
}

// Summarise the comment threads into unresolved items.
func Summarise(ctx context.Context, gc *gerrit.Client, changeID string) (*Summary, error) {
	gcc := &gerrit.ChangesClient{Client: gc}

	ch, err := gcc.GetChange(ctx, changeID, "MESSAGES", "DETAILED_LABELS", "CURRENT_REVISION", "CURRENT_COMMIT", "DETAILED_ACCOUNTS")
	if err != nil {
		return nil, fmt.Errorf("could not get change: %w", err)
	}

	// Extract commit message
	commitMessage := ""
	if len(ch.Revisions) == 1 {
		for _, x := range ch.Revisions {
			commitMessage = x.Commit.Message
		}
	}

	reviewers := ch.Reviewers["REVIEWER"]
	cced := ch.Reviewers["CC"]

	var activeReviewers []gerrit.AccountInfo
	activeReviewersDedup := make(map[string]bool)
	for _, m := range ch.Messages {
		if activeReviewersDedup[m.Author.Username] {
			continue
		}
		activeReviewers = append(activeReviewers, *m.Author)
		activeReviewersDedup[m.Author.Username] = true
	}

	if ch.UnresolvedCommentCount == 0 {
		return &Summary{
			ChangeID:            strconv.Itoa(ch.Number),
			Project:             ch.Project,
			Branch:              ch.Branch,
			Subject:             ch.Subject,
			LatestCommitMessage: commitMessage,
			Created:             ch.Created.Time(),
			Updated:             ch.Updated.Time(),
			Submitted:           ch.Submitted.Time(),
			Comments:            ch.TotalCommentCount,
			UnresolvedComments:  ch.UnresolvedCommentCount,
			AllReviewers:        reviewers,
			ActiveReviewers:     activeReviewers,
			CCed:                cced,
		}, nil
	}

	comments, err := gcc.ListChangeComments(ctx, changeID)
	if err != nil {
		return nil, fmt.Errorf("could not list change comments: %w", err)
	}

	threads := make(map[string]gerrit.CommentInfo)   // Last processed Comment ID -> Latest comment in a thread
	authors := make(map[string][]gerrit.AccountInfo) // Last processed Comment ID -> Authors from the thread
	for path, cs := range comments {
		for _, c := range cs {
			if c.Path == "" {
				c.Path = path
			}
			var as []gerrit.AccountInfo
			// Remove the comment that `c` is replying to...
			if c.InReplyTo != "" {
				delete(threads, c.InReplyTo)

				as = authors[c.InReplyTo]
				delete(authors, c.InReplyTo)
			}

			as = append(as, c.Author)
			authors[c.ID] = as

			// Only record unresolved comments...
			if c.Unresolved == true {
				threads[c.ID] = c
			}
		}
	}

	ucs := make([]gerrit.CommentInfo, 0, len(threads))
	for _, c := range threads {
		ucs = append(ucs, c)
	}

	sort.Slice(ucs, func(i, j int) bool {
		return ucs[i].Updated.Time().Before(ucs[j].Updated.Time())
	})

	for k, as := range authors {
		dedup := make(map[string]struct{})
		out := make([]gerrit.AccountInfo, 0, len(as))
		for _, a := range as {
			if _, ok := dedup[a.Username]; ok {
				continue
			}
			dedup[a.Username] = struct{}{}
			out = append(out, a)
		}
		authors[k] = as
	}

	s := &Summary{
		ChangeID:            strconv.Itoa(ch.Number),
		Project:             ch.Project,
		Branch:              ch.Branch,
		Subject:             ch.Subject,
		LatestCommitMessage: commitMessage,
		Created:             ch.Created.Time(),
		Updated:             ch.Updated.Time(),
		Submitted:           ch.Submitted.Time(),
		Comments:            ch.TotalCommentCount,
		UnresolvedComments:  ch.UnresolvedCommentCount,
		AllReviewers:        reviewers,
		ActiveReviewers:     activeReviewers,
		CCed:                cced,
		Threads:             make([]Thread, 0, len(ucs)),
	}

	for _, uc := range ucs {
		s.Threads = append(s.Threads, Thread{
			s:           s,
			Path:        uc.Path,
			Line:        uc.Line,
			PatchSet:    uc.PatchSet,
			Authors:     authors[uc.ID],
			Message:     uc.Message,
			LastComment: uc,
		})
	}
	return s, nil
}
