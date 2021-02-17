package multiplexer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Router отделён от оператора, т.к. это позволяет нам проще оперировать тем, как мы обрабатываем наши запросы
// имплементации роутера, далеко не важен весь http.Request или способ формирования ответа
type Router interface {
	HandleRequest(w http.ResponseWriter, req *http.Request)
}

type defaultRouter struct {
	Operator
}

func NewDefaultRouter(op Operator) Router {
	return &defaultRouter{Operator: op}
}

const maxURLsInRequest = 20

func (r *defaultRouter) HandleRequest(w http.ResponseWriter, req *http.Request) {
	// method check
	if req.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errPathSupportsPostOnly)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		writeError(w, http.StatusBadRequest, errRequestBodyIsNotJSON)
		return
	}

	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, &errRequestFailedReadBody{err: err})
		return
	}

	body := new(RequestBody)
	err = json.Unmarshal(data, body)
	if err != nil {
		writeError(w, http.StatusBadRequest, &errRequestFailedReadBody{err: err})
		return
	}

	// для проверки тела лучше, возможно, использовать что-то вроде github.com/go-playground/validator
	// но нам не нужно проверять миллиард параметров, так что обойдемся
	if len(*body) > maxURLsInRequest {
		writeError(w, http.StatusBadRequest, errTooManyURLsPerRequest)
		return
	}
	var wrongURLs = make([]string, 0)
	for _, item := range *body {
		_, err := url.Parse(item)
		if err != nil {
			wrongURLs = append(wrongURLs, item)
		}
	}
	if len(wrongURLs) > 0 {
		writeError(w, http.StatusBadRequest, &errParameterIsNotURL{indexes: wrongURLs})
		return
	}

	params := &RequestParameters{
		Ctx:  req.Context(),
		Body: body,
	}

	resp, err := r.Operator.Request(params)
	if err != nil {
		statusCode := 0
		switch err.(type) {
		case *errFetchingEndpoint:
			statusCode = http.StatusBadGateway
		default:
			statusCode = http.StatusInternalServerError
		}

		writeError(w, statusCode, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
