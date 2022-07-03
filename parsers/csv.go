package parsers

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/opendefinition/tuoda/database"
)

type ColumnHeaders struct {
	Preparse bool     `json:"preparse"`
	LinePos  int      `json:"line_pos"`
	Names    []string `json:"column_names"`
	SkipCols []int    `json:"skip_cols`
}

func (ch *ColumnHeaders) StandardizeColumns() {
	for index, value := range ch.Names {
		ch.Names[index] = strings.ToLower(strings.ReplaceAll(value, ".", "_"))
	}
}

type CsvDefinition struct {
	ParserType     string        `json:"parser_type"`
	CommentChar    string        `json:"comment_char"`
	Delimiter      string        `json:"delimiter"`
	ColumnsHeaders ColumnHeaders `json:"column_headers"`
}

func (cd *CsvDefinition) Parse(database database.ArangoDB, collection string, logPath string) {
	logfile, err := os.Open(logPath)

	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()

	csvreader := csv.NewReader(logfile)

	// Check if delimiter char has been set by user
	delimiter, size := utf8.DecodeRuneInString(cd.Delimiter)

	if size > 0 {
		csvreader.Comma = delimiter
	}

	// Enable lazy quotes
	csvreader.LazyQuotes = true

	// Handle commented lines
	if len(cd.CommentChar) > 0 {
		commentchar, char_size := utf8.DecodeRuneInString(cd.CommentChar)

		if char_size > 0 {
			csvreader.Comment = commentchar
		}
	}

	line_counter := 0

	// Test if we need to preparse column headers
	if cd.ColumnsHeaders.Preparse == true && len(cd.ColumnsHeaders.Names) == 0 && cd.ColumnsHeaders.LinePos >= 0 {
		cd.PreParseHeaderColumns(logPath)
	}

	// Make sure that column headers follows a certain naming convention
	cd.ColumnsHeaders.StandardizeColumns()

	for {
		line_counter++

		line, read_err := csvreader.Read()

		if read_err == io.EOF {
			break
		}

		if read_err != nil {
			log.Fatal(read_err)
		}

		if read_err != nil {
			log.Fatal(read_err)
		}

		// Parse headers
		if line_counter == cd.ColumnsHeaders.LinePos {
			cd.ParseHeaderColumns(line)
			continue
		}

		// Parse log line
		entry, parseerr := cd.ParseLogLine(line)

		if parseerr != nil {
			fmt.Printf("Error: %v\n", parseerr)
		} else {
			database.InsertLogItem(collection, entry)
		}
	}
}

func (cd *CsvDefinition) ParseHeaderColumns(line []string) {
	for index, value := range line {
		illegal := false

		for _, skip := range cd.ColumnsHeaders.SkipCols {
			if (index + 1) == skip {
				illegal = true
			}
		}

		if illegal == false {
			cd.ColumnsHeaders.Names = append(cd.ColumnsHeaders.Names, value)
		}
	}
}

func (cd *CsvDefinition) PreParseHeaderColumns(filePath string) {
	logfile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(logfile)
	linecounter := 0

	for scanner.Scan() {
		linecounter++
		if linecounter == cd.ColumnsHeaders.LinePos {
			line := strings.Split(scanner.Text(), cd.Delimiter)
			cd.ParseHeaderColumns(line)
		}
	}

	logfile.Close()
}

func (cd *CsvDefinition) ParseLogLine(logline []string) (map[string]interface{}, error) {
	logentry := make(map[string]interface{})

	length_headers := len(cd.ColumnsHeaders.Names)
	length_columns := len(logline)

	if length_headers == 0 {
		return logentry, errors.New("No CSV headers defined")
	}

	if length_headers != length_columns && length_headers < length_columns {
		return logentry, errors.New("Column headers does not match logline")
	}

	if length_columns == 0 {
		return logentry, errors.New("Encountered empty log line")
	}

	for index, value := range logline {
		logentry[cd.ColumnsHeaders.Names[index]] = value
	}

	// Document id for log entry
	guid := uuid.New()
	logentry["_key"] = guid

	// Document hash
	sha256id := sha256.New()
	sha256id.Write([]byte(fmt.Sprintf("%v", logentry)))
	logentry["docsum"] = fmt.Sprintf("%x", sha256id.Sum(nil))

	return logentry, nil
}
