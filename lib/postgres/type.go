package postgres

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"regexp"
	"strconv"
	"strings"
)

var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	var str string
	switch src := src.(type) {
	case []byte:
		str = string(src)
	case string:
		str = src
	case nil:
		*s = nil
		return nil
	}

	str = quoteEscapeRegex.ReplaceAllString(str, `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	str = str[1 : len(str)-1]

	if len(str) == 0 {
		*s = []string{}
		return nil
	}

	csvReader := csv.NewReader(strings.NewReader(str))
	slice, err := csvReader.Read()
	if err != nil {
		return err
	}
	*s = slice

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(strconv.Quote(val))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return string(buffer.Bytes()), nil
}
