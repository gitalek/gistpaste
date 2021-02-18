package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_secureHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})


	type args struct {
		next http.Handler
		w *httptest.ResponseRecorder
		r *http.Request
	}

	type result struct {
		headerFrameOptions  string
		headerXssProtection string
		body                string
		statusCode          int
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "first",
			args: args{next: next, w: w, r: r},
			want: result{
				headerXssProtection: "1; mode=block",
				headerFrameOptions: "deny",
				body: "OK",
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secureHeaders(next).ServeHTTP(tt.args.w, tt.args.r)
			rs := tt.args.w.Result()
			defer rs.Body.Close()

			if frameOptions := rs.Header.Get("X-Frame-Options"); frameOptions != tt.want.headerFrameOptions {
				t.Errorf("X-Frame-Options: want %s; got %s", tt.want.headerFrameOptions, frameOptions)
			}

			if xssProtection := rs.Header.Get("X-XSS-Protection"); xssProtection != tt.want.headerXssProtection {
				t.Errorf("X-XSS-Protection: want %s; got %s", tt.want.headerXssProtection, xssProtection)
			}

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
