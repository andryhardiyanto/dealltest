package postgres

import (
	"fmt"

	"github.com/pkg/errors"
)

func Pairs(kv []interface{}) (map[string]interface{}, error) {
	if len(kv)%2 == 1 {
		return nil, errors.New(fmt.Sprintf("kv got the odd number of input pairs %d", len(kv)))
	}
	arg := map[string]interface{}{}
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		arg[key] = kv[i+1]
	}
	return arg, nil
}

func PairsHook(kv []interface{}, ids map[string]int64, hook string) (map[string]interface{}, error) {
	if len(kv)%2 == 1 {
		return nil, errors.New(fmt.Sprintf("kv got the odd number of input pairs %d", len(kv)))
	}
	arg := map[string]interface{}{}
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		val := kv[i+1]
		s, ok := val.(string)
		if ok && len(s) > 11 && s[:11] == hook {
			val = ids[s[11:]]
		}
		arg[key] = val
	}
	return arg, nil
}

func Filter(sl []string, m map[string]string) (result []string) {
	for _, v := range sl {
		_, ok := m[v]
		if !ok {
			result = append(result, v)
		}
	}
	return
}
