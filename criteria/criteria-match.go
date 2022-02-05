package criteria

import (
	"errors"
	"fmt"
	"strings"
)

type criteriaMatch struct {
	parent InnerCriteria

	field  string
	substr string

	valFunc ValueTransformer
}

// FieldFulltext is a special field name used for fulltext search
const FieldFulltext string = "--(--fulltext--)--"

// FieldDefault is a special field name used to set the default match constructor
const FieldDefault string = "--(--default--)--"

var ErrCriteriaConstructorFailed = errors.New("criteria constructor error")

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

	return false
}

func (c criteriaMatch) String() string {
	if c.field == FieldFulltext {
		return fmt.Sprintf("ANYWHERE CONTAINS \"%s\"", c.substr)
	}

	return fmt.Sprintf("%s CONTAINS \"%s\"", c.field, c.substr)
}
