package criteria

type Criteria interface {
	Matches([]string) bool
}
