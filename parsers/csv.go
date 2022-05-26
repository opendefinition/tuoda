package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func (cd *CsvDefinition) Parse(logPath string) {
	logFile, err := os.Open(logPath)

	if err != nil {
		fmt.Println(err)
	}

	logScanner := bufio.NewScanner(logFile)

	counter := 1

	for logScanner.Scan() {
		// Extract column names
		if counter == int(cd.ColumnsHeaders.LinePos) {
			if len(cd.ColumnsHeaders.Names) == 0 {
				cd.ParseHeaderColumns(logScanner.Text())
			}
		}

		// Parse logline
		if counter >= (cd.Data.StartsAtLine) {
			logentry := cd.ParseLogLine(logScanner.Text())
			fmt.Println(logentry)
		}

		counter++
	}

	logFile.Close()
}

func (cd *CsvDefinition) ParseHeaderColumns(line string) {
	// Remove unwanted characters
	line = strings.ReplaceAll(line, ".", "_")

	// Removing unwanted columns
	for index, name := range strings.Split(line, cd.Delimiter) {
		illegal := false
		for _, skip := range cd.ColumnsHeaders.SkipCols {
			if (index + 1) == skip {
				illegal = true
			}
		}

		if illegal == false {
			cd.ColumnsHeaders.Names = append(cd.ColumnsHeaders.Names, name)
		}
	}
}

func (cd *CsvDefinition) ParseLogLine(logline string) map[string]interface{} {
	data := strings.Split(logline, cd.Delimiter)
	logentry := make(map[string]interface{})

	for index, value := range data {
		logentry[cd.ColumnsHeaders.Names[index]] = value
	}

	return logentry
}
