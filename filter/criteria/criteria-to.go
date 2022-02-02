package criteria

import (
	"fmt"

	"github.com/emersion/go-imap"
)

type criteriaTo struct {
	to string
}

func (c criteriaTo) Matches(msg *imap.Message) bool {
	return addressesContains(msg.Envelope.To, c.to)
}

func NewTo(to string) criteriaTo {
	return criteriaTo{
		to,
	}
}

func (c criteriaTo) String() string {
	return fmt.Sprintf("to:(%s)", c.to)
}
