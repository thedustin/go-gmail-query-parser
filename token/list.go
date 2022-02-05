package token

import "strings"

// List is a list of tokens, used to add some handy methods.
type List []Token

// String transforms this list into a textual representation using the values of the tokens.
func (tl List) String() string {
	strs := make([]string, len(tl))

	for i, t := range tl {
		strs[i] = t.queryString()
	}

	return strings.Join(strs, "")
}

// Describe transforms this list into a string representation using the type of the tokens.
// Useful for debugging to see the actual types insted of the textual representation.
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

// Valid validates the list to be syntactic fine.
func (tl List) Validate() error {
	for i := 0; i < len(tl); i++ {
		t := tl[i]
		var next *Token

		if i+1 < len(tl) {
			next = &tl[i+1]
		}

		if err := t.Validate(next); err != nil {
			return err
		}
	}

	return nil
}

func kindInList(needle kind, kindList []kind) bool {
	for _, kind := range kindList {
		if needle == kind {
			return true
		}
	}

	return false
}
