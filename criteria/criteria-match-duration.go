package criteria

import (
	"fmt"
	"time"
)

type durationOperator string

const (
	durationOlderThan durationOperator = ">"
	durationNewerThan durationOperator = "<"
)

type criteriaMatchDuration struct {
	parent InnerCriteria

	field    string
	duration time.Duration
	operator durationOperator

	valFunc ValueTransformer
}

func NewMatchDuration(field, duration string, valFunc ValueTransformer, operator durationOperator) *criteriaMatchDuration {
	d, _ := time.ParseDuration(duration) // @todo: Error handling @todo: allow days

	return &criteriaMatchDuration{
		field:    field,
		duration: d,
		operator: operator,
		valFunc:  valFunc,
	}
}

var NewerThanMatchConstructor = CriteriaMatchConstructor(func(field, substr string, valFunc ValueTransformer) InnerCriteria {
	return NewMatchDuration(field, substr, valFunc, durationNewerThan)
})

var OlderThanMatchConstructor = CriteriaMatchConstructor(func(field, substr string, valFunc ValueTransformer) InnerCriteria {
	return NewMatchDuration(field, substr, valFunc, durationOlderThan)
})

func (c criteriaMatchDuration) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaMatchDuration) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c criteriaMatchDuration) Matches(v interface{}) bool {
	now := time.Now()
	vals := c.valFunc(c.field, v)

	for _, v := range vals {
		date, err := time.Parse(time.RFC3339, v)

		if err != nil {
			continue
		}

		var m bool

		switch c.operator {
		case durationNewerThan:
			m = date.Add(c.duration).After(now)
		case durationOlderThan:
			m = date.Add(c.duration).Before(now)
		}

		if m {
			return true
		}
	}

	return false
}

func (c criteriaMatchDuration) String() string {
	var s string

	if c.operator == durationNewerThan {
		s = "NEWER THAN"
	} else {
		s = "OLDER THAN"
	}

	return fmt.Sprintf("%s \"%s\"", s, c.duration)
}
