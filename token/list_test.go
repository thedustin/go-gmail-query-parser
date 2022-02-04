package token_test

import (
	"errors"
	"testing"

	"github.com/thedustin/go-gmail-query-parser/token"
)

type testcase struct {
	Name string
	List token.List
	Err  error
}

func TestListValidation(t *testing.T) {
	ts := []testcase{
		{
			Name: "Empty Query",
			List: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name: "Simple filter",
			List: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "from"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "example.org"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name: "Normal Tuesday",
			List: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "from"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "@example.org"),
				token.NewToken(token.GroupStart, "("),
				token.NewToken(token.Field, "subject"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "Werbung für Treppenlifte"),
				token.NewToken(token.Or, "OR"),
				token.NewToken(token.Negate, "-"),
				token.NewToken(token.Field, "older_than"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "7d"),
				token.NewToken(token.GroupEnd, ")"),
				token.NewToken(token.Fulltext, "from"),
				token.NewToken(token.Fulltext, "Lorem"),
				token.NewToken(token.Fulltext, "ipsum"),
				token.NewToken(token.End, "$"),
			},
		},
	}

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.List.Validate()

			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}

			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}

			if tc.Err != nil && err != nil && !errors.Is(err, tc.Err) {
				t.Errorf("Unexpected error %v", err)
			}
		})
	}
}
