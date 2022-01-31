package v1_test

import (
	gosrouter "github.com/fikryfahrezy/gosrouter/v1"
	"net/http"
	"reflect"
	"testing"
)

func getOne(w http.ResponseWriter, r *http.Request) {}

func getTwo(w http.ResponseWriter, r *http.Request) {}

func getThree(w http.ResponseWriter, r *http.Request) {}

func getFour(w http.ResponseWriter, r *http.Request) {}

func postOne(w http.ResponseWriter, r *http.Request) {}

func postTwo(w http.ResponseWriter, r *http.Request) {}

func postThree(w http.ResponseWriter, r *http.Request) {}

func postFour(w http.ResponseWriter, r *http.Request) {}

func putOne(w http.ResponseWriter, r *http.Request) {}

func putTwo(w http.ResponseWriter, r *http.Request) {}

func deleteOne(w http.ResponseWriter, r *http.Request) {}

func deleteTwo(w http.ResponseWriter, r *http.Request) {}

func TestGetRoute(t *testing.T) {
	gosrouter.Routes = make(map[string]gosrouter.RouteChild)

	cases := []struct {
		testName, regUrl, reqUrl, mtd string
		regFn                         func(url string, fn http.HandlerFunc)
		fn                            func(http.ResponseWriter, *http.Request)
	}{
		{
			"post handler",
			"/",
			"/",
			"POST",
			gosrouter.HandlerPOST,
			postOne,
		},
		{
			"get handler",
			"/",
			"/",
			"GET",
			gosrouter.HandlerGET,
			getOne,
		},
		{
			"post handler with param dynamic id",
			"/:id",
			"/1",
			"POST",
			gosrouter.HandlerPOST,
			postTwo,
		},
		{
			"get handler with param dynamic id",
			"/:id",
			"/1",
			"GET",
			gosrouter.HandlerGET,
			getTwo,
		},
		{
			"put handler with param dynamic id",
			"/:id",
			"/1",
			"PUT",
			gosrouter.HandlerPUT,
			putOne,
		},
		{
			"delete handler with param dynamic id",
			"/:id",
			"/1",
			"DELETE",
			gosrouter.HandlerDELETE,
			deleteOne,
		},
		{
			"post handler with static param",
			"/one",
			"/one",
			"POST",
			gosrouter.HandlerPOST,
			postThree,
		},
		{
			"get handler with static param",
			"/one",
			"/one",
			"GET",
			gosrouter.HandlerGET,
			getThree,
		},
		{
			"post handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"POST",
			gosrouter.HandlerPOST,
			postFour,
		},
		{
			"get handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"GET",
			gosrouter.HandlerGET,
			getFour,
		},
		{
			"put handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"PUT",
			gosrouter.HandlerPUT,
			putTwo,
		},
		{
			"delete handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"DELETE",
			gosrouter.HandlerDELETE,
			deleteTwo,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		t.Run(v.testName, func(t *testing.T) {
			if rt := gosrouter.GetRoute(v.reqUrl, v.mtd); reflect.ValueOf(rt).Pointer() != reflect.ValueOf(v.fn).Pointer() {
				t.Fatal(v.regUrl)
			}
		})
	}
}

func TestDynamicRoute(t *testing.T) {
	gosrouter.Routes = make(map[string]gosrouter.RouteChild)

	cases := []struct {
		testName, regUrl, reqUrl, paramName, param, mtd string
		regFn                                           func(url string, fn http.HandlerFunc)
		fn                                              func(http.ResponseWriter, *http.Request)
	}{
		{
			"get handler with dynamic param at depth 0",
			"/:id",
			"/1",
			"id",
			"1",
			"GET",
			gosrouter.HandlerGET,
			getOne,
		},
		{
			"get handler with dynamic param at depth 1",
			"/v1/:x",
			"/v1/11",
			"x",
			"11",
			"GET",
			gosrouter.HandlerGET,
			getOne,
		},
		{
			"get handler with dynamic param at depth 2",
			"/v2/v3/:xy",
			"/v2/v3/xyz",
			"xy",
			"xyz",
			"GET",
			gosrouter.HandlerGET,
			getOne,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		t.Run(v.testName, func(t *testing.T) {
			p := gosrouter.ReqParams(v.reqUrl, v.paramName)

			if p != v.param {
				t.Fatal(v.regUrl)
			}
		})
	}
}

func BenchmarkDynamicRoute(b *testing.B) {
	gosrouter.Routes = make(map[string]gosrouter.RouteChild)
	gosrouter.HandlerGET("/:id", getOne)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		_ = gosrouter.ReqParams("/1", "id")
	}
}

func BenchmarkDynamicRoute5(b *testing.B) {
	gosrouter.Routes = make(map[string]gosrouter.RouteChild)

	regUrl := "/:a/:b/:c/:d/:e"
	reqUrl := "/1/2/3/4/5"
	gosrouter.HandlerGET(regUrl, getOne)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		_ = gosrouter.ReqParams(reqUrl, "a")
		_ = gosrouter.ReqParams(reqUrl, "b")
		_ = gosrouter.ReqParams(reqUrl, "c")
		_ = gosrouter.ReqParams(reqUrl, "d")
		_ = gosrouter.ReqParams(reqUrl, "e")
	}
}

func BenchmarkDynamicRoute20(b *testing.B) {
	gosrouter.Routes = make(map[string]gosrouter.RouteChild)

	regUrl := "/:a/:b/:c/:d/:e/:f/:g/:h/:i/:j/:k/:l/:m/:n/:o/:p/:q/:r/:s/:t"
	reqUrl := "/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20"
	gosrouter.HandlerGET(regUrl, getOne)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		_ = gosrouter.ReqParams(reqUrl, "a")
		_ = gosrouter.ReqParams(reqUrl, "b")
		_ = gosrouter.ReqParams(reqUrl, "c")
		_ = gosrouter.ReqParams(reqUrl, "d")
		_ = gosrouter.ReqParams(reqUrl, "e")
		_ = gosrouter.ReqParams(reqUrl, "f")
		_ = gosrouter.ReqParams(reqUrl, "g")
		_ = gosrouter.ReqParams(reqUrl, "h")
		_ = gosrouter.ReqParams(reqUrl, "i")
		_ = gosrouter.ReqParams(reqUrl, "j")
		_ = gosrouter.ReqParams(reqUrl, "k")
		_ = gosrouter.ReqParams(reqUrl, "l")
		_ = gosrouter.ReqParams(reqUrl, "m")
		_ = gosrouter.ReqParams(reqUrl, "n")
		_ = gosrouter.ReqParams(reqUrl, "o")
		_ = gosrouter.ReqParams(reqUrl, "p")
		_ = gosrouter.ReqParams(reqUrl, "q")
		_ = gosrouter.ReqParams(reqUrl, "r")
		_ = gosrouter.ReqParams(reqUrl, "s")
		_ = gosrouter.ReqParams(reqUrl, "t")
	}
}
