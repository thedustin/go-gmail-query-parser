package token

import "strings"

type List []Token

func (tl List) String() string {
	strs := make([]string, len(tl))

	for i, t := range tl {
		strs[i] = t.queryString()
	}

	return strings.Join(strs, "")
}

func (tl List) Describe() string {
	strs := make([]string, len(tl))

	for i, t := range tl {
		v := string(t.kind)

		if t.isTrailingSpaceToken() {
			v = v + " "
		}

		strs[i] = v
	}

	return strings.Join(strs, "")
}

func KindInList(t kind, vs []kind) bool {
	for _, vT := range vs {
		if t == vT {
			return true
		}
	}

	return false
}
