package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	parser "github.com/thedustin/go-gmail-query-parser"
	"github.com/thedustin/go-gmail-query-parser/criteria"
)

func main() {
	f := criteria.ValueTransformer(func(field string, v interface{}) []string {
		return []string{
			"from:john.doe@example.org",
			"subject:Werbung für Treppenlifte",
			"from Lorem ipsum",
		}
	})

	p := parser.NewParser(f, parser.FlagDefault|parser.FlagOptimize)

	crit, err := p.Parse(strings.Join(os.Args[1:], ""))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", strings.Join(os.Args[1:], ""))
	fmt.Printf("%s (%#v)\n", crit, crit)

	fmt.Println("Does it match?", crit.Matches(nil))
}
