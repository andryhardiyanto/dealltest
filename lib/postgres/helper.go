package postgres

import (
	"fmt"
	"strings"

	"github.com/AndryHardiyanto/dealltest/lib/log"
)

const (
	qInsert = "insert"
	qUpdate = "update"
	qDelete = "delete"
	qSelect = "select"

	qResult = "q-result---"
)

func debugQuery(query string, args map[string]interface{}) {
	d := log.Debug()
	for k, v := range args {
		d.Str(k, fmt.Sprintf("%v", v))
	}
	replacer := strings.NewReplacer("\n", "", "\t", "")
	d.Msg(replacer.Replace(query))
}

func queryType(query string) string {
	query = strings.TrimSpace(query)
	query = query[:6]
	if strings.EqualFold(query, qInsert) {
		return qInsert
	} else if strings.EqualFold(query, qUpdate) {
		return qUpdate
	} else if strings.EqualFold(query, qDelete) {
		return qDelete
	} else {
		return qSelect
	}
}
