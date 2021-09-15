`gosrouter` stands for Go Simple Router, which was made to sufficient author need to get dynamic params at [author Back-End project](https://github.com/fikryfahrezy/gobookshelf) when author learning Go for the first time.

## Install

`go get -u github.com/fikryfahrezy/gosrouter`

## Example

```go
package main

import (
	"net/http"

	"github.com/fikryfahrezy/gosrouter"
)

func post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi from post"))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi from get"))
}

func delete(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		w.Write([]byte("id required"))
		return
	}

	w.Write([]byte("hi from delete"))
}

func put(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	x := p("x")

	if x == "" {
		w.Write([]byte("x required"))
		return
	}

	w.Write([]byte("hi from put"))
}

func patch(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		w.Write([]byte("id required"))
		return
	}

	w.Write([]byte("hi from patch"))
}

func main() {
	gosrouter.HandlerPOST("/", post)
	gosrouter.HandlerGET("/", get)
	gosrouter.HandlerDELETE("/:id", delete)
	gosrouter.HandlerPUT("/:x", put)
	gosrouter.HandlerPATCH("/x/:id", patch)

	for r := range gosrouter.Routes {
		http.HandleFunc(r, gosrouter.MakeHandler)
	}

	http.ListenAndServe(":3000", nil)
}
```

## License

Free to use for any projects.
