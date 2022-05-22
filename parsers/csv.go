package parsers

import "fmt"

type Test struct {
	Jalla string
}

type ColumnHeaders struct {
	LinePos  uint     `json:"line_pos"`
	Names    []string `json:"column_names"`
	SkipCols []int    `json:"skip_cols`
}

type DataRow struct {
	StartsAtLine int `json:"starts_at_line"`
}

type CsvDefinition struct {
	ParserType     string        `json:"parser_type"`
	Delimiter      string        `json:"delimiter"`
	ColumnsHeaders ColumnHeaders `json:"column_headers"`
	Data           DataRow       `json:"data"`
}

func (cd *CsvDefinition) Parse() {
	fmt.Println("Parsing CSV")
}
