package criteria

import (
	"fmt"
	"regexp"
	"strconv"
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

func NewMatchDuration(field, duration string, valFunc ValueTransformer, operator durationOperator) (*criteriaMatchDuration, error) {
	d, err := extendedDurationFormat(duration)

	if err != nil {
		return nil, err
	}

	return &criteriaMatchDuration{
		field:    field,
		duration: d,
		operator: operator,
		valFunc:  valFunc,
	}, nil
}

var NewerThanMatchConstructor = CriteriaMatchConstructor(func(field, substr string, valFunc ValueTransformer) (InnerCriteria, error) {
	return NewMatchDuration(field, substr, valFunc, durationNewerThan)
})

var OlderThanMatchConstructor = CriteriaMatchConstructor(func(field, substr string, valFunc ValueTransformer) (InnerCriteria, error) {
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

var r = regexp.MustCompile("(([0-9]+)y)?(([0-9]+)m)?(([0-9]+)w)?(([0-9]+)d)?")

func extendedDurationFormat(dur string) (time.Duration, error) {
	dur = r.ReplaceAllStringFunc(dur, func(s string) string {
		if s == "" {
			return s
		}

		l := len(s)
		v, err := strconv.Atoi(s[:l-1])

		if err != nil {
			return s
		}

		switch s[l-1] {
		case 'y':
			return fmt.Sprintf("%dh", (v * 365 * 24))
		case 'm':
			return fmt.Sprintf("%dh", (v * 30 * 24))
		case 'w':
			return fmt.Sprintf("%dh", (v * 7 * 24))
		case 'd':
			return fmt.Sprintf("%dh", (v * 24))
		}

		return s
	})

	return time.ParseDuration(dur)
}
