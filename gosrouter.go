package gosrouter

import (
	"fmt"
	"net/http"
	"strings"
)

type RouteChild struct {
	Depth   int
	Route   string
	Dynamic bool
	Fn      map[string]func(http.ResponseWriter, *http.Request)
	Child   *RouteChild
}

var routeMethods = map[string]string{
	"GET":    http.MethodGet,
	"POST":   http.MethodPost,
	"PUT":    http.MethodPut,
	"DELETE": http.MethodDelete,
	"PATCH":  http.MethodPatch,
}

var Routes = make(map[string]RouteChild)

func routeChild(r *RouteChild, i, m int, mtd string, s []string, fn func(http.ResponseWriter, *http.Request)) RouteChild {
	var nr RouteChild

	if r != nil {
		nr = *r
	}

	rt := s[i]
	nr.Depth = i
	nr.Route = fmt.Sprintf("/%s", rt)

	if strings.HasPrefix(rt, ":") {
		nr.Dynamic = true
	} else {
		nr.Dynamic = false
	}

	if i == m-1 {
		if f := nr.Fn; f == nil {
			nr.Fn = make(map[string]func(http.ResponseWriter, *http.Request))
		}

		nr.Fn[mtd] = fn
		return nr
	}

	i++

	if nc := routeChild(nr.Child, i, m, mtd, s, fn); nc.Route != "" {
		nr.Child = &nc
	}

	return nr
}

func registerHandler(mtd, url string, fn func(http.ResponseWriter, *http.Request)) {
	if strings.Contains(url, ":") {
		s := strings.Split(url, "/")[1:]
		l := len(s)

		if l == 0 {
			return
		}

		fe := s[0]
		r := "/"

		if !strings.HasPrefix(fe, ":") {
			r += fe
		} else {
			// Ref: How to prepend int to slice
			// https://stackoverflow.com/questions/53737435/how-to-prepend-int-to-slice
			s = append(s, "")
			copy(s[1:], s)
			s[0] = ""
			l = len(s)
		}

		o := Routes[r]
		Routes[r] = routeChild(&o, 0, l, mtd, s, fn)
	} else {
		o := Routes[url]

		if o.Route == "" {
			o = RouteChild{Depth: 0, Route: url}
		}

		if f := o.Fn; f == nil {
			o.Fn = make(map[string]func(http.ResponseWriter, *http.Request))
		}

		o.Fn[mtd] = fn

		if Routes[url].Route == "" {
			Routes[url] = o
		}
	}
}

func HandlerPOST(url string, fn http.HandlerFunc) {
	registerHandler("POST", url, fn)
}

func HandlerGET(url string, fn http.HandlerFunc) {
	registerHandler("GET", url, fn)
}

func HandlerPUT(url string, fn http.HandlerFunc) {
	registerHandler("PUT", url, fn)
}

func HandlerDELETE(url string, fn http.HandlerFunc) {
	registerHandler("DELETE", url, fn)
}

func HandlerPATCH(url string, fn http.HandlerFunc) {
	registerHandler("PATCH", url, fn)
}

func getRoute(url, mtd string) func(http.ResponseWriter, *http.Request) {
	if r := Routes[url]; r.Route == url && r.Fn[mtd] != nil {
		return r.Fn[mtd]
	}

	s := strings.Split(strings.Replace(url, "/", " /", -1), " ")[1:]

	if len(s) == 1 {
		if h := Routes[s[0]].Fn[mtd]; h != nil {
			return h
		}
	}

	var l RouteChild

	for i, v := range s {
		if i == 0 {
			if r, rc := Routes[v], Routes["/"].Child; r.Route == "" && rc != nil {
				l = *rc
			} else {
				l = r
			}
		}

		if f := l.Fn[mtd]; f != nil && i == len(s)-1 {
			return f
		}

		if l.Child != nil {
			l = *l.Child
		}
	}

	return nil
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	m := routeMethods[strings.ToUpper(r.Method)]

	if m == "" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	rt := getRoute(r.URL.Path, m)

	if rt == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	rt(w, r)
}

func ReqParams(u, p string) string {
	s := strings.Split(strings.Replace(u, "/", " /", -1), " ")[1:]

	var l RouteChild
	isSls := false

	for i, v := range s {
		if i == 0 {
			if r := Routes[v]; r.Route == "" {
				lt := Routes["/"]
				l = *lt.Child
				isSls = true
			} else {
				l = r
			}
		}

		if l.Dynamic && strings.Split(l.Route, "/:")[1] == p {
			if isSls {
				return s[l.Depth-1][1:]
			}
			return s[l.Depth][1:]
		}

		if l.Child != nil {
			l = *l.Child
		}
	}

	return ""
}
