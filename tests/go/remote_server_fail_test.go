//go:build integration
// +build integration

package previewer_test

import (
	"net/http"
	"testing"
)

func TestBadServer(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/fail_server.com/some_image.jpg")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadGateway {
		t.Errorf("on a non-existing server, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusBadGateway)
	}
}

func TestBadImage(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/nginx/bad_image.jpg")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("on a non-existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusNotFound)
	}
}
