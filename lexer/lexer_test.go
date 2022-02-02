package lexer_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/thedustin/go-gmail-query-parser/lexer"
	"github.com/thedustin/go-gmail-query-parser/token"
)

type testcase struct {
	Name   string
	Source string
	Err    error
	Result token.List
}

func TestParser(t *testing.T) {
	ts := []testcase{
		{
			Name: "Empty Query",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "Simple filter",
			Source: "from:example.org",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "from"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "example.org"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "Complexe filter value",
			Source: "subject:(Werbung für Treppenlifte)",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "subject"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "Werbung für Treppenlifte"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "Negate filter",
			Source: "-older_than:7d",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Negate, "-"),
				token.NewToken(token.Field, "older_than"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "7d"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "OR filter",
			Source: "older_than:7d OR larger:2M",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "older_than"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "7d"),
				token.NewToken(token.Or, "OR"),
				token.NewToken(token.Field, "larger"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "2M"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "AND filter",
			Source: "older_than:7d AND larger:2M",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Field, "older_than"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "7d"),
				token.NewToken(token.Field, "larger"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "2M"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "Normal Tuesday",
			Source: "from:(@example.org) (subject:(Werbung für Treppenlifte) OR -older_than:7d) from Lorem ipsum",
			Result: token.List{
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
		{
			Name:   "No syntax check at this point",
			Source: "--older_than:7d",
			Result: token.List{
				token.NewToken(token.Start, "^"),
				token.NewToken(token.Negate, "-"),
				token.NewToken(token.Negate, "-"),
				token.NewToken(token.Field, "older_than"),
				token.NewToken(token.Equal, ":"),
				token.NewToken(token.FieldValue, "7d"),
				token.NewToken(token.End, "$"),
			},
		},
		{
			Name:   "Lets break the groups",
			Source: "foo (older_than:7d bar",
			Err:    lexer.ErrGroupNotClosed,
		},
	}

	l := lexer.NewLexer()

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			err := l.Parse(tc.Source)

			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}

			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}

			if tc.Err != nil && err != nil && !errors.Is(err, tc.Err) {
				t.Errorf("Unexpected error %v", err)
			}

			if !reflect.DeepEqual(tc.Result, l.Result()) {
				t.Errorf("Result does not match, expected\n\t%#v\nbut got\n\t%#v", tc.Result, l.Result())
			}

			t.Logf("%s", tc.Source)
			t.Logf("%s", l.Result())
			t.Logf("%s", l.Result().Describe())
		})
	}
}
