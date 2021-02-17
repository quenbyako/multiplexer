package multiplexer

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}

func writeJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(obj)
	if err != nil {
		writeError(w, http.StatusBadRequest, errServerJSONEncodingError)
		return
	}

	w.Write(data)
}

func isDataBinary(in []byte) bool {
	str := string(in) // т.к. байты проверить не сможем

	for _, char := range str {
		if !unicode.IsPrint(char) && !strings.ContainsRune("\r\n\t", char) {
			return true
		}
	}

	return false
}

func chunkSlice(in []string, chunkCount int) [][]string {
	if chunkCount <= 0 {
		panic("invalid chunk count")
	}
	if chunkCount == 1 {
		return [][]string{in}
	}

	if len(in) <= chunkCount {
		res := make([][]string, len(in))
		for i, item := range in {
			res[i] = []string{item}
		}
		return res
	}

	res := make([][]string, 0, chunkCount)

	chunkSize := (len(in) + chunkCount - 1) / chunkCount
	for i := 0; i < len(in); i += chunkSize {
		end := i + chunkSize

		if end > len(in) {
			end = len(in)
		}

		res = append(res, in[i:end])
	}

	return res
}
