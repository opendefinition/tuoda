package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

const TYPE_STRING = "STRING"
const TYPE_FLOAT = "FLOAT"
const TYPE_INTEGER = "INTEGER"

func DetermineType(value interface{}) (string, string) {
	regexes := map[string]string{
		TYPE_FLOAT:   "^\\d+[,.]\\d+$",
		TYPE_INTEGER: "^\\d+$",
	}

	unpacked := reflect.ValueOf(value).String()
	for key, regex := range regexes {
		match, err := regexp.MatchString(regex, unpacked)

		if err != nil {
			fmt.Println("Unable to cast type")
		}

		if match == true {
			return key, unpacked
		}
	}

	return TYPE_STRING, unpacked
}

func StringToFLoat(value string) float64 {
	s, err := strconv.ParseFloat(value, 64)

	if err != nil {
		fmt.Println("Unable to cast to float64")
	}

	return s
}

func StringToInteger(value string) int64 {
	s, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		fmt.Println("Unable to cast to integer64")
	}

	return s
}
