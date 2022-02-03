package criteria

type criteriaAnd struct {
	criterias []Criteria
}

func (c criteriaAnd) Matches(vals []string) bool {
	for _, subCrit := range c.criterias {
		if !subCrit.Matches(vals) {
			return false
		}
	}

	return true
}

func NewAnd(criterias ...Criteria) criteriaAnd {
	return criteriaAnd{criterias}
}

func (c *criteriaAnd) Add(crit Criteria) {
	c.criterias = append(c.criterias, crit)
}
