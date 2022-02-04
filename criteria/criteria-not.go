package criteria

type criteriaNot struct {
	parent   InnerCriteria
	criteria InnerCriteria
}

func NewNot(criteria InnerCriteria) *criteriaNot {
	return &criteriaNot{
		criteria: criteria,
	}
}

func (c criteriaNot) Matches(v interface{}) bool {
	return !c.criteria.Matches(v)
}

func (c criteriaNot) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaNot) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c criteriaNot) String() string {
	return "NOT " + c.criteria.String() + ""
}
