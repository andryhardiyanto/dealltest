package errors

import "net/http"

var (
	ErrorMessages  map[string]*Message
	DefaultMessage = &Message{
		Language: "EN",
		Message:  "Server error occurred, Please try again later",
		Type:     TypeInternalServerError,
		Code:     http.StatusInternalServerError,
	}
)

type Message struct {
	Language string `json:"-"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Code     int    `json:"-"`
}

type errorMessage struct {
	errorMsgs map[string]map[string]*Message
}

type ErrorMessage interface {
	Translate(request *http.Request, err error) *Message
}

func RegisterErrorMessage() ErrorMessage {
	return &errorMessage{
		errorMsgs: Msg,
	}
}

func (em *errorMessage) Translate(request *http.Request, err error) *Message {
	errMsg, ok := err.(*Error)
	if ok {
		return em.mappingTranslateMessage(request, errMsg.Type)
	}

	return em.mappingTranslateMessage(request, TypeInternalServerError)
}

func (em *errorMessage) mappingTranslateMessage(request *http.Request, errType string) *Message {
	language := "ID"
	if request != nil && request.Header.Get("language") != "" {
		language = request.Header.Get("language")
	}

	msg, ok := em.errorMsgs[language]
	if !ok {
		return DefaultMessage
	}

	return msg[errType]
}
