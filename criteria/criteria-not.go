package criteria

type criteriaNot struct {
	criteria Criteria
}

func (c criteriaNot) Matches(vals []string) bool {
	return !c.criteria.Matches(vals)
}

func NewNot(criteria Criteria) criteriaNot {
	return criteriaNot{
		criteria: criteria,
	}
}
