`gosrouter` stands for Go Simple Router, which was made to sufficient author need to get dynamic params
at [author Back-End project](https://github.com/fikryfahrezy/gobookshelf) when author learning Go for the first time.

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
	id := gosrouter.ReqParams(r.URL.Path, "id")

	if id == "" {
		w.Write([]byte("id required"))
		return
	}

	w.Write([]byte("hi from delete"))
}

func put(w http.ResponseWriter, r *http.Request) {
	x := gosrouter.ReqParams(r.URL.Path, "x")

	if x == "" {
		w.Write([]byte("x required"))
		return
	}

	w.Write([]byte("hi from put"))
}

func patch(w http.ResponseWriter, r *http.Request) {
	id := gosrouter.ReqParams(r.URL.Path, "id")

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

## Benchmark
```
go test -bench=. -gcflags -m ./...

goos: windows
goarch: amd64
pkg: github.com/fikryfahrezy/gosrouter
cpu: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
BenchmarkDynamicRoute-8          4967024               249.6 ns/op            67 B/op          3 allocs/op
BenchmarkDynamicRoute5-8          394708              2937 ns/op            1040 B/op         25 allocs/op
BenchmarkDynamicRoute20-8          33561             35093 ns/op           15360 B/op        250 allocs/op
```

## License

MIT
