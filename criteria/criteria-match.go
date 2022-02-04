package criteria

import (
	"fmt"
	"strings"
)

type criteriaMatch struct {
	parent InnerCriteria

	field  string
	substr string

	valFunc ValueTransformer
}

const FieldFulltext string = "(--fulltext)"

func NewMatch(field, substr string, valFunc ValueTransformer) *criteriaMatch {
	return &criteriaMatch{
		field:   field,
		substr:  substr,
		valFunc: valFunc,
	}
}

func (c criteriaMatch) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaMatch) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c criteriaMatch) Matches(v interface{}) bool {
	vals := c.valFunc(c.field, v)

	for _, v := range vals {
		if strings.Contains(v, c.substr) {
			return true
		}
	}

	fmt.Println("[", "does not match", "]", c.field, c.substr)

	return false
}

func (c criteriaMatch) String() string {
	if c.field == FieldFulltext {
		return fmt.Sprintf("CONTAINS \"%s\"", c.substr)
	}

	return fmt.Sprintf("%s == \"%s\"", c.field, c.substr)
}
