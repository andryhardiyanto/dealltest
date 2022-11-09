package gin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/gin-gonic/gin"
)

type request struct {
	ctx *gin.Context
}

func NewRequest(ctx *gin.Context) *request {
	return &request{
		ctx: ctx,
	}
}

func (r *request) GetRequest() map[string]interface{} {
	var (
		dataBody interface{}
	)

	dataQueryParam := r.GetQuery()
	r.GetBody(&dataBody)
	dataParam := r.GetParams()
	request := map[string]interface{}{
		"param":  dataParam,
		"body":   dataBody,
		"query":  dataQueryParam,
		"form":   r.ctx.Request.Form,
		"header": r.ctx.Request.Header,
	}

	return request
}

func (r *request) GetBody(data interface{}) error {
	bytesBody, _ := ioutil.ReadAll(r.ctx.Request.Body)
	r.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bytesBody))
	err := json.Unmarshal(bytesBody, &data)
	if err != nil {
		return errors.NewWrapError(err, "error unmarshal request").SetType(errors.TypeUnprocessableEntity)
	}

	return nil
}
func (r *request) GetParams() gin.Params {
	return r.ctx.Params
}
func (r *request) GetJsonQuery(data interface{}) error {
	query := r.GetQuery()
	bytesBody, _ := json.Marshal(query)

	err := json.Unmarshal(bytesBody, &data)
	if err != nil {
		return errors.NewWrapError(err, "error unmarshal request").SetType(errors.TypeUnprocessableEntity)
	}

	return nil
}
func (r *request) GetQuery() map[string]interface{} {
	data := make(map[string]interface{})
	queryParam := strings.Split(r.ctx.Request.URL.RawQuery, "\u0026")

	if len(queryParam) > 0 && queryParam[0] != "" {
		queryParamMap := make(map[string][]string)
		for _, qp := range queryParam {
			query := strings.Split(qp, "=")
			if len(query) == 2 {
				queryParamMap[query[0]] = append(queryParamMap[query[0]], query[1])
			} else if len(query) == 1 {
				queryParamMap[query[0]] = append(queryParamMap[query[0]], "")
			}
		}
		queryParamMapInterface := make(map[string]interface{})
		for key, qp := range queryParamMap {
			if len(qp) == 1 {
				queryParamMapInterface[key] = qp[0]
			} else if len(qp) > 1 {
				queryParamMapInterface[key] = qp
			}
		}
		data = queryParamMapInterface
	}

	return data
}
