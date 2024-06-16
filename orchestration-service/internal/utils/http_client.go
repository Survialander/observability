package utils

import (
	"crypto/tls"
	"net/http"
)

func GetHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	client := &http.Client{Transport: tr}
	return client
}
