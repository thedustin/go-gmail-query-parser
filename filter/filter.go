package filter

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Filter struct {
	criteria Criteria
	actions  []Action
}

func NewFilter(criteria Criteria, actions ...Action) Filter {
	return Filter{criteria, actions}
}

func (f *Filter) Execute(msg *imap.Message, c *client.Client) error {
	if !f.criteria.Matches(msg) {
		return nil
	}

	for _, a := range f.actions {
		err := a.Perform(msg, c)

		if err != nil {
			return err
		}
	}

	return nil
}
