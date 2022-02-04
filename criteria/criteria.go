package criteria

import "errors"

type Criteria interface {
	Matches(v interface{}) bool
}

type InnerCriteria interface {
	Criteria

	Parent() InnerCriteria
	SetParent(InnerCriteria)
	String() string
}

type GroupCriteria interface {
	Criteria

	Add(InnerCriteria) error
	Replace(old InnerCriteria, new InnerCriteria) error
	Index(InnerCriteria) int
	Length() int

	All() []InnerCriteria
}

var ErrCriteriaAlreadyInGroup = errors.New("criteria already in group")
var ErrCriteriaNotInGroup = errors.New("criteria not in group")
