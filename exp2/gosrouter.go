package exp2

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RoutesChild struct {
	Urls, Params []string
	Child        map[string]http.HandlerFunc
}

type RouteChildV2 map[string]RoutesChild

var routeMethods = map[string]string{
	"GET":    http.MethodGet,
	"POST":   http.MethodPost,
	"PUT":    http.MethodPut,
	"DELETE": http.MethodDelete,
	"PATCH":  http.MethodPatch,
}

var Routes RouteChildV2

func Handler(mtd, url string, handler http.HandlerFunc) {
	urls := strings.Split(url, "/")[1:]
	urlsLen := len(urls)
	if urlsLen == 0 {
		urlsLen = 1
	}

	var s strings.Builder
	s.WriteString("^")
	us := make([]string, urlsLen)
	prs := make([]string, urlsLen)
	if url == "/" {
		s.WriteString("/")
		us[0] = ""
		prs[0] = ""
	} else {
		for i, u := range urls {
			if strings.HasPrefix(u, ":") {
				s.WriteString("/[[:word:]]+")
				us[i] = "%s"
				prs[i] = u
			} else {
				s.WriteString("/" + u)
				us[i] = u
				prs[i] = ""
			}
		}
	}

	s.WriteString("$")
	ss := s.String()

	if _, ok := Routes[ss]; !ok {
		Routes[ss] = RoutesChild{
			Urls:   us,
			Params: prs,
			Child:  make(map[string]http.HandlerFunc),
		}
	}

	Routes[ss].Child[mtd] = handler
}

func ReqParams(url, p string) string {
	urls := strings.Split(url, "/")[1:]

	for k, v := range Routes {
		matched, err := regexp.MatchString(k, url)
		if err != nil || !matched {
			continue
		}

		var s string
		for i, ur := range v.Params {
			if strings.TrimPrefix(ur, ":") == p {
				fmt.Sscanf(urls[i], v.Urls[i], &s)
			}
		}

		if s != "" {
			return s
		}
	}

	return ""
}

func HandlerGET(url string, handler http.HandlerFunc) {
	Handler("GET", url, handler)
}

func HandlerPOST(url string, handler http.HandlerFunc) {
	Handler("POST", url, handler)
}

func HandlerPUT(url string, handler http.HandlerFunc) {
	Handler("PUT", url, handler)
}

func HandlerDELETE(url string, handler http.HandlerFunc) {
	Handler("DELETE", url, handler)
}

func GetRoute(url, mtd string) http.HandlerFunc {
	m := strings.ToUpper(mtd)

	var rc RoutesChild
	rc, ok := Routes["^"+url+"$"]
	if ok && rc.Child != nil {
		route, ok := rc.Child[m]
		if !ok {
			return nil
		}

		return route
	}

	for k, v := range Routes {
		matched, err := regexp.MatchString(k, url)
		if err != nil || !matched {
			continue
		}
		rc = v
	}

	if rc.Child != nil {
		route, ok := rc.Child[m]
		if !ok {
			return nil
		}

		return route
	}

	return nil
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	m := routeMethods[strings.ToUpper(r.Method)]
	if m == "" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if route := GetRoute(r.URL.Path, m); route != nil {
		route(w, r)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
