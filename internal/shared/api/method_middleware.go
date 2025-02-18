package api

import "net/http"

type Middleware func(http.Handler) http.Handler

func Method(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method != method {
			methodNotAllowed(response)

			return
		}

		handler(response, request)
	}
}
