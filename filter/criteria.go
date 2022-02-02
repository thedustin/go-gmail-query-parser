package filter

import "github.com/emersion/go-imap"

type Criteria interface {
	Matches(msg *imap.Message) bool
	String() string
}
