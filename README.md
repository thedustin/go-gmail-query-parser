# [Go Gmail Query Parser](https://github.com/thedustin/go-gmail-query-parser)

by [Dustin Breuer](https://github.com/thedustin)

The Go Gmail Query Parser is a simple parser for [Gmail-like filter queries/search operators][gmail-search-operators] like this: `from:@example.org (subject:(Lorem ipsum) OR older_than:7d)`. The goal of this project was to create a lightweight parser in plain Go that I can use in my other project to filter and sort emails from different inboxes, and to learn how a simple parser works.

As you can customize the supported fields and values, the parser is not limitied to email sorting.

# Prerequisites

Before installing this project you need:

* 🐀 Go (at least 1.17)


## Installation

To build this project on your local machine, just run `go build`.


## Usage

```golang
// see "Support custom values" to understand what "f" is
p := parser.NewParser(f, parser.FlagDefault)

// create an expression from our filter string
expr, err := p.Parse(filterQueryStr)

if err != nil {
    log.Fatalln(err)
}

// filter some values with our expression
for _, v := range values {
    if expr.Matches(v) {
        v.DoSomeWork()
    }
}
```

## Customize parser

The parser is really customizable and allows altering supported fields and values. Some words explained:

| Word | Description |
|------|-------------|
| Field | A Field is a string used as keyword to filter for a specific property. Some of the built-in fields are e.g. `from`, `to`, `subject`, `older_than`, `before`,... |
| Criteria | A Criteria is an object that does the actual work.

### Support custom values

The parser requires a `type ValueTransformer func(field string, v interface{}) []string` as argument. This function will be called to fetch the values for a specific field from the object to filter.

```golang
f := criteria.ValueTransformer(func(field string, v interface{}) []string {
    user, ok := v.(*app.User)

    if !ok {
        return []string{} // or panic
    }

    switch field: {
    case "username":
        return []string{user.Username}
    case "email":
        return []string{user.Email}
    case "last_login":
        return []string{user.LastLogin.Format(time.RFC3339)}
    case criteria.FieldFulltext:
        return []string{
            user.Username,
            user.Email,
            user.ProfileText,
        }
    }

    return []string{}
})

p := parser.NewParser(f, parser.FlagDefault|parser.FlagOptimize)
```

### Add new fields

Adding new fields is quite easy. We just need the name of the field and a constructor function that will be called to create a criteria later.

```golang
p := parser.NewParser(f, parser.FlagDefault|parser.FlagOptimize)

// Add the fields and use the default criteria for these fields
// These calls will fail if we already added the fields
p.AddField("username", DefaultCriteriaMatchConstructor)
p.AddField("email", nil)

// Add the field and use a "newer than" constructor to use a duration as value
// If there is already a field with that name, we just overwrite it
p.setField("last_login_before", NewerThanMatchConstructor)
```

### Remove fields

You can remove specific fields with `func (p *Parser) RemoveField(name string) error`.

To remove all fields use `func (p *Parser) RemoveAllFields(name string)`.


## Contributing to this project

If you have suggestions for improving the prometheus exporter, please
[open an issue or pull request on GitHub](https://github.com/thedustin/crushftp-prometheus-exporter/).


## Additional links

- [Search operators you can use with Gmail][gmail-search-operators]
- [Create rules to filter your emails][gmail-create-filter-rules]

[gmail-search-operators]: https://support.google.com/mail/answer/7190
[gmail-create-filter-rules]: https://support.google.com/mail/answer/6579
