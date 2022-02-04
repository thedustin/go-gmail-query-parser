package criteria

import "strings"

type criteriaAnd struct {
	parent    InnerCriteria
	criterias []InnerCriteria
}

func NewAnd(criterias ...InnerCriteria) *criteriaAnd {
	return &criteriaAnd{
		criterias: criterias,
	}
}

func (c criteriaAnd) Matches(v interface{}) bool {
	for _, subCrit := range c.criterias {
		if !subCrit.Matches(v) {
			return false
		}
	}

	return true
}

func (c criteriaAnd) Parent() InnerCriteria {
	return c.parent
}

func (c *criteriaAnd) SetParent(p InnerCriteria) {
	c.parent = p
}

func (c *criteriaAnd) Add(crit InnerCriteria) error {
	if i := c.Index(crit); i != -1 {
		return ErrCriteriaAlreadyInGroup
	}

	crit.SetParent(c)

	c.criterias = append(c.criterias, crit)

	return nil
}

func (c criteriaAnd) Index(needle InnerCriteria) int {
	for i, haystack := range c.criterias {
		if haystack == needle {
			return i
		}
	}

	return -1
}

func (c *criteriaAnd) Replace(old InnerCriteria, new InnerCriteria) error {
	i := c.Index(old)

	if i == -1 {
		return ErrCriteriaNotInGroup
	}

	old.SetParent(nil)
	new.SetParent(c)

	c.criterias[i] = new

	return nil
}

func (c criteriaAnd) Length() int {
	return len(c.criterias)
}

func (c criteriaAnd) All() []InnerCriteria {
	return c.criterias
}

func (c criteriaAnd) String() string {
	crits := make([]string, len(c.criterias))

	for i := 0; i < len(c.criterias); i++ {
		crits[i] = c.criterias[i].String()
	}

	return "(" + strings.Join(crits, " ") + ")"
}
