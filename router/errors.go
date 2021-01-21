package router

import (
	"fmt"
	"net/http"
)

type httpErr struct {
	err error
	msg string
}

func (he httpErr) Error() string {
	return he.err.Error()
}

func ErrJSON(w http.ResponseWriter, err error, code int) {

	v, ok := err.(httpErr)
	if !ok {
		fmt.Printf("An error occurred: %v", err)
		http.Error(w, err.Error(), code)
		return
	}
	fmt.Printf("An error occurred: %v", v.err)
	JSON(w, map[string]string{
		"error": v.msg,
	}, code)

}
