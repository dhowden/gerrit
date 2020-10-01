package thread

import (
	"context"
	"fmt"
	"sort"

	"github.com/dhowden/gerrit"
)

// Summary of a change.
type Summary struct {
	ChangeID string
	Project  string

	Comments           int
	UnresolvedComments int

	Threads []Thread
}

// Thread of comments.
type Thread struct {
	Path     string
	Line     int
	PatchSet int
	Authors  []gerrit.AccountInfo
	Message  string

	LastComment gerrit.CommentInfo
}

// Summarise the comment threads into unresolved items.
func Summarise(ctx context.Context, gc *gerrit.Client, changeID string) (*Summary, error) {
	gcc := &gerrit.ChangesClient{Client: gc}

	ch, err := gcc.GetChange(ctx, changeID)
	if err != nil {
		return nil, fmt.Errorf("could not get change: %w", err)
	}

	if ch.UnresolvedCommentCount == 0 {
		return &Summary{Comments: ch.TotalCommentCount}, nil
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

	ts := make([]Thread, 0, len(ucs))
	for _, uc := range ucs {
		ts = append(ts, Thread{
			Path:        uc.Path,
			Line:        uc.Line,
			PatchSet:    uc.PatchSet,
			Authors:     authors[uc.ID],
			Message:     uc.Message,
			LastComment: uc,
		})
	}
	return &Summary{
		ChangeID:           changeID,
		Project:            ch.Project,
		Comments:           ch.TotalCommentCount,
		UnresolvedComments: ch.UnresolvedCommentCount,
		Threads:            ts,
	}, nil
}
