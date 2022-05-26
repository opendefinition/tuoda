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
			cd.ParseHeaderColumns(logScanner.Text())
		}

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
	cd.ColumnsHeaders = strings.Split(line, cd.Delimiter)

	for _, index := range cd.ColumnsHeaders.SkipCols {
		index -= 1
		cd.ColumnsHeaders = append(cd.ColumnsHeaders[:index], cd.ColumnsHeaders[index+1:]...)
	}
}

func (cd *CsvDefinition) ParseLogLine(logline string) []string {
	line := strings.Split(logline, cd.Delimiter)
	fmt.Println(line)

	return line
}
