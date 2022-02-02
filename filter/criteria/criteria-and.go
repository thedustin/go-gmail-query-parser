package criteria

import (
	"strings"

	"github.com/emersion/go-imap"
	"github.com/thedustin/go-email-curator/filter"
)

type criteriaAnd struct {
	criterias []filter.Criteria
}

func (c criteriaAnd) Matches(msg *imap.Message) bool {
	for _, subCrit := range c.criterias {
		if !subCrit.Matches(msg) {
			return false
		}
	}

	return true
}

func NewAnd(criterias ...filter.Criteria) criteriaAnd {
	return criteriaAnd{criterias}
}

func (c criteriaAnd) String() string {
	crits := make([]string, len(c.criterias))

	for i := 0; i < len(c.criterias); i++ {
		crits[i] = c.criterias[i].String()
	}

	return "(" + strings.Join(crits, " ") + ")"
}
