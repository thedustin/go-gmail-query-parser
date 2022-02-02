package criteria

import (
	"fmt"

	"github.com/emersion/go-imap"
)

type criteriaFrom struct {
	from string
}

func (c criteriaFrom) Matches(msg *imap.Message) bool {
	return addressesContains(msg.Envelope.From, c.from)
}

func NewFrom(from string) criteriaFrom {
	return criteriaFrom{
		from,
	}
}

func (c criteriaFrom) String() string {
	return fmt.Sprintf("from:(%s)", c.from)
}
