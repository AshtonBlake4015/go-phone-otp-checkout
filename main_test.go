package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyCaptchaRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/captcha/verify" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"widget_record_id": "widget-1", "token": "token-1"}
		if len(body) != len(want) || body["widget_record_id"] != want["widget_record_id"] || body["token"] != want["token"] {
			t.Fatalf("body = %#v, want %#v", body, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{"success":true}}`))
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "test-key", HTTP: server.Client()}
	result, err := c.VerifyCaptcha("widget-1", "token-1")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != `{"success":true}` {
		t.Fatalf("result = %s", result)
	}
}
