package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const signatureHeader = "X-TRAQ-Signature"

var client = &http.Client{}

func CalcHMACSHA1(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func PostWebhook(url, message, secret string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("creating webhook request: %w", err)
	}
	req.Header.Set(signatureHeader, CalcHMACSHA1(message, secret))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		var body []byte
		if resp != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		return fmt.Errorf("posting webhook (body: %s): %w", string(body), err)
	}
	return nil
}
