package criteria

import "strings"

type criteriaOr struct {
	parent    InnerCriteria
	criterias []InnerCriteria
}

func NewOr(criterias ...InnerCriteria) *criteriaOr {
	return &criteriaOr{
		criterias: criterias,
	}
}

func (c criteriaOr) Matches(v interface{}) bool {
	for _, subCrit := range c.criterias {
		if subCrit.Matches(v) {
			return true
		}
	}

	return false
}

func (c criteriaOr) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaOr) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c *criteriaOr) Add(crit InnerCriteria) error {
	if i := c.Index(crit); i != -1 {
		return ErrCriteriaAlreadyInGroup
	}

	crit.SetParent(c)

	c.criterias = append(c.criterias, crit)

	return nil
}

func (c criteriaOr) Index(needle InnerCriteria) int {
	for i, haystack := range c.criterias {
		if haystack == needle {
			return i
		}
	}

	return -1
}

func (c *criteriaOr) Replace(old InnerCriteria, new InnerCriteria) error {
	i := c.Index(old)

	if i == -1 {
		return ErrCriteriaNotInGroup
	}

	old.SetParent(nil)
	new.SetParent(c)

	c.criterias[i] = new

	return nil
}

func (c criteriaOr) Length() int {
	return len(c.criterias)
}

func (c criteriaOr) All() []InnerCriteria {
	return c.criterias
}

func (c criteriaOr) String() string {
	crits := make([]string, len(c.criterias))

	for i := 0; i < len(c.criterias); i++ {
		crits[i] = c.criterias[i].String()
	}

	return "(" + strings.Join(crits, " OR ") + ")"
}
