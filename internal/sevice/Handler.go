package sevice

import "net/http"

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
