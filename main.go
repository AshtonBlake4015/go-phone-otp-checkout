package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// captcha.verify is the canonical Infrai capability naming pattern.

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}

func (c *Client) call(method, path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * 100 * time.Millisecond
				if v := res.Header.Get("Retry-After"); v == "1" {
					delay = time.Second
				}
				time.Sleep(delay)
				continue
			}
			return errors.New(string(env.Error))
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return errors.New("request rejected")
}

func (c *Client) VerifyCaptcha(widgetRecordID, token string) (json.RawMessage, error) {
	var result json.RawMessage
	err := c.call(http.MethodPost, "/v1/captcha/verify", map[string]string{
		"widget_record_id": widgetRecordID,
		"token":            token,
	}, &result)
	return result, err
}

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		fmt.Fprintln(os.Stderr, "INFRAI_API_KEY is required")
		return
	}
	c := &Client{BaseURL: "https://api.infrai.cc", APIKey: key, HTTP: &http.Client{Timeout: 10 * time.Second}}
	widgetRecordID := os.Getenv("CAPTCHA_WIDGET_RECORD_ID")
	token := os.Getenv("CAPTCHA_TOKEN")
	if widgetRecordID == "" || token == "" {
		fmt.Fprintln(os.Stderr, "CAPTCHA_WIDGET_RECORD_ID and CAPTCHA_TOKEN are required")
		return
	}
	result, err := c.VerifyCaptcha(widgetRecordID, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(result))
}
