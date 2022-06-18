package parsers

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/opendefinition/tuoda/config"
	"github.com/opendefinition/tuoda/database"
)

type ColumnHeaders struct {
	LinePos  int      `json:"line_pos"`
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

func (cd *CsvDefinition) Parse(config config.Configuration, logPath string) {
	logfile, err := os.Open(logPath)

	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()

	csvreader := csv.NewReader(logfile)

	if cd.Delimiter == "\t" {
		csvreader.Comma = '\t'
	}

	// Prepare Arango database
	Arango := database.ArangoDBClient(
		config.ArangoDB.Address,
		config.ArangoDB.Database,
		config.ArangoDB.Username,
		config.ArangoDB.Password,
	)

	// Ask where to store the log
	var collection_name string
	fmt.Print("Name of collection: ")
	fmt.Scanln(&collection_name)
	fmt.Println("")

	line_counter := 1

	for {
		line, err := csvreader.Read()

		if err == io.EOF {
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(line_counter)
		// Obtain CSV headers from indicated line
		if cd.ColumnsHeaders.LinePos >= 0 && line_counter == cd.ColumnsHeaders.LinePos {
			cd.ParseHeaderColumns(line)
			line_counter++
			continue
		}

		if line_counter >= cd.Data.StartsAtLine {
			logentry := cd.ParseLogLine(line)
			Arango.InsertLogItem(collection_name, logentry)
		}

		line_counter++
	}
}

func (cd *CsvDefinition) ParseHeaderColumns(line []string) {
	for index, value := range line {
		column_name := strings.ReplaceAll(value, ".", "_")

		illegal := false

		for _, skip := range cd.ColumnsHeaders.SkipCols {
			if (index + 1) == skip {
				illegal = true
			}
		}

		if illegal == false {
			cd.ColumnsHeaders.Names = append(cd.ColumnsHeaders.Names, column_name)
		}
	}
}

func (cd *CsvDefinition) ParseLogLine(logline []string) map[string]interface{} {
	logentry := make(map[string]interface{})

	for index, value := range logline {
		logentry[cd.ColumnsHeaders.Names[index]] = value
	}

	// Generate document id for log entry
	sha256id := sha256.New()
	sha256id.Write([]byte(fmt.Sprintf("%v", logentry)))
	logentry["_key"] = fmt.Sprintf("%x", sha256id.Sum(nil))

	return logentry
}
