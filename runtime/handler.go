package runtime

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Evaler evaluates script content sent via HTTP.
type Evaler interface {
	Eval([]byte) (any, error)
}

// ServeMux is a mux abstraction.
type ServeMux interface {
	Handle(string, http.Handler)
}

// AddHandler registers the evaluation handler with a ServeMux.
func AddHandler(mux ServeMux, evaler Evaler) {
	mux.Handle("/eval", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := evaler.Eval(body); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		fmt.Fprintf(w, "OK")
	}))
}
