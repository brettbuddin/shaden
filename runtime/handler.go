package runtime

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Evaler interface {
	Eval([]byte) (interface{}, error)
}

type ServeMux interface {
	Handle(string, http.Handler)
}

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
