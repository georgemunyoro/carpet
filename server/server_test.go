package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var fh = fileHandler{root: http.Dir("../testing/test_site/"), projectDir: "../testing/test_site/"}
var handler = http.HandlerFunc(fh.ServeHTTP)

func TestServeHTTPIndexHTML(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	temp, err := ioutil.ReadFile("../testing/test_site/index.html")
	if err != nil {
		t.Fatal(err)
	}
	TestSiteIndexHTML := string(temp) + eventSourceScript
	if strings.TrimSpace(rr.Body.String()) != TestSiteIndexHTML {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), TestSiteIndexHTML)
	}
}

func TestServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/example_post.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	temp, err := ioutil.ReadFile("../testing/test_site/posts/example_post.html")
	if err != nil {
		t.Fatal(err)
	}
	HTML := string(temp) + eventSourceScript
	if strings.TrimSpace(rr.Body.String()) != HTML {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), HTML)
	}
}

func TestServeHTTP404(t *testing.T) {
	req, err := http.NewRequest("GET", "/non-existent/page", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
