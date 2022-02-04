package criteria

type criteriaNoop struct {
	parent InnerCriteria
}

func NewNoop() *criteriaNoop {
	return &criteriaNoop{}
}

func (c criteriaNoop) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaNoop) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c criteriaNoop) Matches(v interface{}) bool {
	return true
}

func (c criteriaNoop) String() string {
	return "NOOP"
}
