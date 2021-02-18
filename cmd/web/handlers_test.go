package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ping(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	type result struct {
		body       string
		statusCode int
	}

	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "first",
			args: args{w: w, r: r},
			want: result{body: "OK", statusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ping(tt.args.w, tt.args.r)
			rs := tt.args.w.Result()
			defer rs.Body.Close()
			if rs.StatusCode != tt.want.statusCode {
				t.Errorf("status code: want %d; got %d", tt.want.statusCode, rs.StatusCode)
			}
			body, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			if b := string(body); b != tt.want.body {
				t.Errorf("body: want %s; got %s", tt.want.body, b)
			}
		})
	}
}

func TestEndToEndPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	type args struct {
		route string
		ts    *testServer
	}

	type result struct {
		body       string
		statusCode int
	}

	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "first",
			args: args{route: "/ping", ts: ts},
			want: result{body: "OK", statusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := tt.args.ts.get(t, tt.args.route)
			if statusCode != tt.want.statusCode {
				t.Errorf("status code: want %d; got %d", tt.want.statusCode, statusCode)
			}
			if b := string(body); b != tt.want.body {
				t.Errorf("body: want %s; got %s", tt.want.body, b)
			}
		})
	}
}
