package v2

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

var Routes RouteChildV2

func Handler(mtd, url string, handler http.HandlerFunc) {
	urls := strings.Split(url, "/")[1:]
	urlsLen := len(urls)
	if urlsLen == 0 {
		urlsLen = 1
	}

	var s strings.Builder
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

func URLParam(url, p string) string {
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
	var rc RoutesChild
	for k, v := range Routes {
		matched, err := regexp.MatchString(k, url)
		if err != nil || !matched {
			continue
		}
		rc = v
	}

	if rc.Child != nil {
		route, ok := rc.Child[strings.ToUpper(mtd)]
		if !ok {
			return nil
		}

		return route
	}

	return nil
}

func MakeHandler(w http.ResponseWriter, rq *http.Request) {
	url := rq.URL.Path

	if route := GetRoute(url, rq.Method); route != nil {
		route(w, rq)
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}
