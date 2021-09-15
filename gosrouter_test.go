package gosrouter

import (
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
	Routes = make(map[string]RouteChild)

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
			HandlerPOST,
			postOne,
		},
		{
			"get handler",
			"/",
			"/",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"post handler with param dynamic id",
			"/:id",
			"/1",
			"POST",
			HandlerPOST,
			postTwo,
		},
		{
			"get handler with param dynamic id",
			"/:id",
			"/1",
			"GET",
			HandlerGET,
			getTwo,
		},
		{
			"put handler with param dynamic id",
			"/:id",
			"/1",
			"PUT",
			HandlerPUT,
			putOne,
		},
		{
			"delete handler with param dynamic id",
			"/:id",
			"/1",
			"DELETE",
			HandlerDELETE,
			deleteOne,
		},
		{
			"post handler with static param",
			"/one",
			"/one",
			"POST",
			HandlerPOST,
			postThree,
		},
		{
			"get handler with static param",
			"/one",
			"/one",
			"GET",
			HandlerGET,
			getThree,
		},
		{
			"post handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"POST",
			HandlerPOST,
			postFour,
		},
		{
			"get handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"GET",
			HandlerGET,
			getFour,
		},
		{
			"put handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"PUT",
			HandlerPUT,
			putTwo,
		},
		{
			"delete handler with static param and dynamic param id",
			"/one/:id",
			"/one/1",
			"DELETE",
			HandlerDELETE,
			deleteTwo,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		t.Run(v.testName, func(t *testing.T) {
			if rt := getRoute(v.reqUrl, v.mtd); reflect.ValueOf(rt).Pointer() != reflect.ValueOf(v.fn).Pointer() {
				t.Fatal(v.regUrl)
			}
		})
	}
}

func TestDynamicRoute(t *testing.T) {
	Routes = make(map[string]RouteChild)

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
			HandlerGET,
			getOne,
		},
		{
			"get handler with dynamic param at depth 1",
			"/v1/:x",
			"/v1/11",
			"x",
			"11",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"get handler with dynamic param at depth 2",
			"/v2/v3/:xy",
			"/v2/v3/xyz",
			"xy",
			"xyz",
			"GET",
			HandlerGET,
			getOne,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		t.Run(v.testName, func(t *testing.T) {
			p := ReqParams(v.reqUrl)

			if p(v.paramName) != v.param {
				t.Fatal(v.regUrl)
			}
		})
	}
}
