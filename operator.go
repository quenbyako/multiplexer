package multiplexer

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type RequestParameters struct {
	// в случае таймаута со стороны клиента, мы сможем красиво положить реквесты
	Ctx  context.Context
	Body *RequestBody
}

type RequestBody []string

type RequestResponse []string

type Operator interface {
	Request(*RequestParameters) (*RequestResponse, error)
}

type operatorImplementation struct{}

func NewDefaultOperator() Operator {
	return &operatorImplementation{}
}

const maxRequests = 4

func (*operatorImplementation) Request(req *RequestParameters) (*RequestResponse, error) {
	splittedAddresses := chunkSlice(*req.Body, maxRequests)
	resultParted := make([][]string, len(splittedAddresses))

	ctx, cancel := context.WithCancel(req.Ctx)
	defer cancel()

	var wg sync.WaitGroup

	var gotError error
	var errMutex sync.Mutex // мютекс возьмем, что бы у горутин не было гонки, это бы было совсем обидно

	for i, jobs := range splittedAddresses {
		i, jobs := i, jobs
		println("worker started", i)
		wg.Add(1)
		go func() {
			defer println("worker ended", i)
			defer wg.Done()

			res, err := FetchAllEndpoints(ctx, jobs)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				errMutex.Lock()
				if err != nil {
					gotError = err
				}
				errMutex.Unlock()

				cancel() // отменяем все другие горутины, поскольку уже словили ошибку

				return
			}

			resultParted[i] = res
		}()
	}
	wg.Wait()

	if gotError != nil {
		return nil, gotError
	}

	res := make([]string, 0, len(*req.Body))
	for _, workerResponses := range resultParted {
		res = append(res, workerResponses...)
	}

	response := RequestResponse(res)
	return &response, nil
}

// utils

// FetchAllEndpoints работает в виде воркера и последовательно загружает данные. распараллеливанием занимается
// другая функия
func FetchAllEndpoints(ctx context.Context, urls []string) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(urls) == 0 {
		return []string{}, nil
	}

	res := make([]string, len(urls))

	for i, endpoint := range urls {
		select {
		case <-ctx.Done():
			// если что, бросаем ошибку, т.к. далее нет смысла фетчить данные
			return nil, context.Canceled
		default:
			data, err := FetchData(ctx, endpoint)
			if err != nil {
				return nil, &errFetchingEndpoint{url: endpoint, err: err}
			}

			res[i] = data
		}
	}

	return res, nil
}

const defaultRequestMethod = http.MethodGet

// повторю здесь комментарий из readme: в задании об этом не указано, но возможно что ендпоинт может отдать
// бинарные данные, так что на всякий случай если в ответе есть непечатаемые символы, будем кодировать их в
// base64. По хорошему, надо дополнительно еще проверить mimetype ответа, но это отдельная комплексная задача
func FetchData(ctx context.Context, endpoint string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, defaultRequestMethod, endpoint, nil)
	if errors.Is(err, &url.Error{}) {
		return "", fmt.Errorf("invalid url: '%v'", endpoint)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", &errRequestFailedDesc{err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return "", &errRequestEndedWithBadStatusCode{code: resp.StatusCode, status: resp.Status}
	}

	// FIXME: ioutil не поддерживается начиная с go1.16, но мы предполагаем совместимость с go1.13
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", &errRequestFailedReadBody{err: err}
	}

	var result string

	if isDataBinary(data) {
		result = base64.StdEncoding.EncodeToString(data)
	} else {
		result = string(data)
	}

	return result, nil
}
