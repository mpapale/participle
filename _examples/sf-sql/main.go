package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// Please keep these in lexicographical order
var keywords = []string{
	"DATABASES",
	"DESCRIBE",
	"FROM",
	"LIKE",
	"LIMIT",
	"SCHEMAS",
	"SHOW",
	"STARTS",
	"TABLE",
	"TERSE",
	"VIEW",
	"WITH",
}

type ObjectTypePlural string

const (
	ObjectTypePluralDatabases ObjectTypePlural = "DATABASES"
	ObjectTypePluralSchemas   ObjectTypePlural = "SCHEMAS"
)

func (o *ObjectTypePlural) Capture(values []string) error {
	*o = ObjectTypePlural(values[0])
	return nil
}

type Statement struct {
	Show     *Show     `  @@`
	Describe *Describe `| @@`
}

type Show struct {
	Show             string            `"SHOW"`
	Terse            bool              `(@"TERSE")?`
	ObjectTypePlural *ObjectTypePlural `@("DATABASES"|"SCHEMAS")`
	Like             *string           `("LIKE" @String)?`
	StartsWith       *string           `("STARTS" "WITH" @String)?`
	Limit            *Limit            `("LIMIT" @@)?`
}

type Limit struct {
	Rows *int    `@Number`
	From *string `("FROM" @String)?`
}

type Describe struct {
	Describe string `@("DESCRIBE"|"DESC") @("TABLE"|"VIEW")`
	Name     string `@Ident`
}

var (
	cli struct {
		SQL string `arg:"" required:"" help:"Snowflake SQL to parse."`
	}

	sqlLexer = lexer.Must(lexer.NewSimple([]lexer.Rule{
		{"Keyword", fmt.Sprintf(`(?i)\b(%v)\b`, strings.Join(keywords, "|")), nil},
		// TODO vet with SF
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`, nil},
		// TODO vet with SF
		{"Number", `[-+]?\d*\.?\d+([eE][-+]?\d+)?`, nil},
		// TODO vet with SF
		{`String`, `'[^']*'`, nil},
		{"whitespace", `\s+`, nil},
	}))

	parser = participle.MustBuild(
		&Statement{},
		participle.Lexer(sqlLexer),
		participle.CaseInsensitive("Keyword"),
		participle.Unquote("String"),
	)
)

func main() {
	ctx := kong.Parse(&cli)
	sql := &Statement{}
	err := parser.ParseString("", cli.SQL, sql)
	repr.Println(sql, repr.Indent("  "), repr.OmitEmpty(true))
	ctx.FatalIfErrorf(err)
}
