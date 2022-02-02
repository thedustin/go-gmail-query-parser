package criteria

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
)

type criteriaSubject struct {
	subject string
}

func (c criteriaSubject) Matches(msg *imap.Message) bool {
	return strings.Contains(msg.Envelope.Subject, c.subject)
}

func NewSubject(subject string) criteriaSubject {
	return criteriaSubject{
		subject,
	}
}

func (c criteriaSubject) String() string {
	return fmt.Sprintf("subject:(%s)", c.subject)
}
