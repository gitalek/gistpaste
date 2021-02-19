package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func Test_application_showGist(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	type args struct {
		route string
	}

	type result struct {
		body       []byte
		statusCode int
	}

	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "Valid ID", args: args{route: "/gist/1"}, want: result{body: []byte("An old silent pond..."), statusCode: http.StatusOK}},
		{name: "Non-existent ID", args: args{route: "/gist/2"}, want: result{body: nil, statusCode: http.StatusNotFound}},
		{name: "Negative ID", args: args{route: "/gist/-1"}, want: result{body: nil, statusCode: http.StatusBadRequest}},
		{name: "Decimal ID", args: args{route: "/gist/1.23"}, want: result{body: nil, statusCode: http.StatusBadRequest}},
		{name: "String ID", args: args{route: "/gist/foo"}, want: result{body: nil, statusCode: http.StatusBadRequest}},
		{name: "Empty ID", args: args{route: "/gist/"}, want: result{body: nil, statusCode: http.StatusNotFound}},
		{name: "Trailing slash", args: args{route: "/gist/1/"}, want: result{body: nil, statusCode: http.StatusNotFound}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, tt.args.route)
			if statusCode != tt.want.statusCode {
				t.Errorf("status code: want %d; got %d", tt.want.statusCode, statusCode)
			}
			if !bytes.Contains(body, tt.want.body) {
				t.Errorf("body: want %s; got %s", tt.want.body, body)
			}
			// todo: !
			//if !bytes.Contains(tt.want.body, body) {
			//	t.Errorf("body: want %s; got %s", tt.want.body, body)
			//}
		})
	}
}

func Test_application_signupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	route := "/user/signup"
	_, _, body := ts.get(t, route)
	csrfToken := extractCSRFToken(t, body)

	type formData struct {
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
	}

	type args struct {
		formData formData
	}

	type result struct {
		body       []byte
		statusCode int
	}

	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "Valid submission", args: args{formData: formData{userName: "Bob", userEmail: "bob@example.com", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: nil, statusCode: http.StatusSeeOther}},
		{name: "Empty name", args: args{formData: formData{userName: "", userEmail: "bob@example.com", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field cannot be blank"), statusCode: http.StatusOK}},
		{name: "Empty email", args: args{formData: formData{userName: "Bob", userEmail: "", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field cannot be blank"), statusCode: http.StatusOK}},
		{name: "Empty password", args: args{formData: formData{userName: "Bob", userEmail: "bob@example.com", userPassword: "", csrfToken: csrfToken}}, want: result{body: []byte("This field cannot be blank"), statusCode: http.StatusOK}},
		{name: "Invalid email (incomplete domain)", args: args{formData: formData{userName: "Bob", userEmail: "bob@example.", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field is invalid"), statusCode: http.StatusOK}},
		{name: "Invalid email (missing @)", args: args{formData: formData{userName: "Bob", userEmail: "bobexample.com", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field is invalid"), statusCode: http.StatusOK}},
		{name: "Invalid email (missing local part)", args: args{formData: formData{userName: "Bob", userEmail: "@example.com", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field is invalid"), statusCode: http.StatusOK}},
		{name: "Short password", args: args{formData: formData{userName: "Bob", userEmail: "bob@example.com", userPassword: "pa$$word", csrfToken: csrfToken}}, want: result{body: []byte("This field is too short (minimum is 10 characters)"), statusCode: http.StatusOK}},
		{name: "Duplicate email", args: args{formData: formData{userName: "Bob", userEmail: "error@example.com", userPassword: "validPa$$word", csrfToken: csrfToken}}, want: result{body: []byte("The address is already in use"), statusCode: http.StatusOK}},
		{name: "Invalid CSRF Token", args: args{formData: formData{userName: "", userEmail: "", userPassword: "", csrfToken: "wrongToken"}}, want: result{body: nil, statusCode: http.StatusBadRequest}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.args.formData.userName)
			form.Add("email", tt.args.formData.userEmail)
			form.Add("password", tt.args.formData.userPassword)
			form.Add("csrf_token", tt.args.formData.csrfToken)

			statusCode, _, body := ts.postForm(t, route, form)

			if statusCode != tt.want.statusCode {
				t.Errorf("status code: want %d; got %d", tt.want.statusCode, statusCode)
			}
			if !bytes.Contains(body, tt.want.body) {
				t.Errorf("body: want %s; got %s", tt.want.body, body)
			}
		})
	}
}
