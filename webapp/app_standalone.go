package webapp

import "net/http"

func ListenAndServe(addr string) {
	// FIXME: use absolute path from GOPATH or embed contents
	http.Handle("/", http.FileServer(http.Dir("webapp/static")))
	http.ListenAndServe(addr, nil)
}
