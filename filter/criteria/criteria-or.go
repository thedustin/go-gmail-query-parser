package criteria

import (
	"strings"

	"github.com/emersion/go-imap"
	"github.com/thedustin/go-email-curator/filter"
)

type criteriaOr struct {
	criterias []filter.Criteria
}

func (c criteriaOr) Matches(msg *imap.Message) bool {
	for _, subCrit := range c.criterias {
		if subCrit.Matches(msg) {
			return true
		}
	}

	return false
}

func NewOr(criterias ...filter.Criteria) criteriaOr {
	return criteriaOr{
		criterias: criterias,
	}
}

func (c criteriaOr) String() string {
	crits := make([]string, len(c.criterias))

	for i := 0; i < len(c.criterias); i++ {
		crits[i] = c.criterias[i].String()
	}

	return "(" + strings.Join(crits, " OR ") + ")"
}
