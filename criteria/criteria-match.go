package criteria

import (
	"strings"
)

type criteriaMatch struct {
	substr string
}

func (c criteriaMatch) Matches(vals []string) bool {
	for _, v := range vals {
		if strings.Contains(v, c.substr) {
			return true
		}
	}

	return false
}

func NewMatch(substr string) criteriaMatch {
	return criteriaMatch{
		substr,
	}
}
