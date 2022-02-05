package lexer

import (
	"errors"
	"fmt"
)

type fieldMap map[string]bool

const (
	FieldAfter       = "after"
	FieldBcc         = "bcc"
	FieldBefore      = "before"
	FieldCategory    = "category"
	FieldCc          = "cc"
	FieldDeliveredTo = "deliveredto"
	FieldFilename    = "filename"
	FieldFrom        = "from"
	FieldHas         = "has"
	FieldIn          = "in"
	FieldIs          = "is"
	FieldLabel       = "label"
	FieldLarger      = "larger"
	FieldList        = "list"
	FieldNewerThan   = "newer_than"
	FieldNewer       = "newer"
	FieldOlderThan   = "older_than"
	FieldOlder       = "older"
	FieldRfc822msgid = "rfc822msgid"
	FieldSize        = "size"
	FieldSmaller     = "smaller"
	FieldSubject     = "subject"
	FieldTo          = "to"
)

var defaultFields = fieldMap{
	FieldAfter:       true,
	FieldBcc:         true,
	FieldBefore:      true,
	FieldCategory:    true,
	FieldCc:          true,
	FieldDeliveredTo: true,
	FieldFilename:    true,
	FieldFrom:        true,
	FieldHas:         true,
	FieldIn:          true,
	FieldIs:          true,
	FieldLabel:       true,
	FieldLarger:      true,
	FieldList:        true,
	FieldNewerThan:   true,
	FieldNewer:       true,
	FieldOlderThan:   true,
	FieldOlder:       true,
	FieldRfc822msgid: true,
	FieldSize:        true,
	FieldSmaller:     true,
	FieldSubject:     true,
	FieldTo:          true,
}

var ErrFieldAlreadyDefined = errors.New("field already defined")
var ErrFieldNotDefined = errors.New("field not defined")

func (l *Lexer) AddField(name string) error {
	if _, ok := l.fields[name]; ok {
		return fmt.Errorf("cannot add field %q: %w", name, ErrFieldAlreadyDefined)
	}

	l.fields[name] = true

	return nil
}

func (l *Lexer) SetField(name string, define bool) {
	if _, ok := l.fields[name]; ok && !define {
		delete(l.fields, name)

		return
	}

	l.fields[name] = true
}

func (l *Lexer) RemoveField(name string) error {
	if _, ok := l.fields[name]; ok {
		return fmt.Errorf("cannot remove field %q: %w", name, ErrFieldNotDefined)
	}

	l.fields[name] = true

	return nil
}

func (l *Lexer) RemoveAllFields() {
	l.fields = make(fieldMap)
}
