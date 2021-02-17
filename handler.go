package multiplexer

import (
	"net/http"
)

func SetupHandler(o Router) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/request", o.HandleRequest)

	return r
}
