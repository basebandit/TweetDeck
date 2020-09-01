package api

import (
	"fmt"
	"net/http"
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before")
		h.ServeHTTP(w, r) //call original
		fmt.Println("After")
	})
}
