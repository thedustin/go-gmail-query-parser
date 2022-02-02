package criteria

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
)

type criteriaOlderThan struct {
	duration time.Duration
}

func (c criteriaOlderThan) Matches(msg *imap.Message) bool {
	return msg.Envelope.Date.Before(time.Now().Add(-c.duration))
}

func NewOlderThan(duration time.Duration) criteriaOlderThan {
	return criteriaOlderThan{
		duration,
	}
}

func (c criteriaOlderThan) String() string {
	return fmt.Sprintf("older_than:%s", formatDuration(c.duration))
}
