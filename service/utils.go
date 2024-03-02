package service

import (
	"net/http"
	"time"
)

func DefaultClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
}
