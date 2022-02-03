package criteria

type criteriaOr struct {
	criterias []Criteria
}

func (c criteriaOr) Matches(vals []string) bool {
	for _, subCrit := range c.criterias {
		if subCrit.Matches(vals) {
			return true
		}
	}

	return false
}

func NewOr(criterias ...Criteria) criteriaOr {
	return criteriaOr{
		criterias: criterias,
	}
}

func (c *criteriaOr) Add(crit Criteria) {
	c.criterias = append(c.criterias, crit)
}
