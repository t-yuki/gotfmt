package webapp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/t-yuki/gotfmt/traceback"
)

func init() {
	http.HandleFunc("/process", processHandler)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	var in io.Reader
	if input := r.PostFormValue("traceback"); input != "" {
		in = bytes.NewBufferString(input)
	} else if r.Body != nil {
		in = r.Body
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	trace, err := traceback.ParseTraceback(in, ioutil.Discard)
	if err != nil {
		log.Printf("parse err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	trace.Stacks = traceback.TrimSourcePrefix(trace.Stacks)
	b, err := json.Marshal(trace)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
