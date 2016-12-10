package util

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)

// Convert raw string to cookie.
func StringToCookie(cookie string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", cookie)
	tsav := &http.Request{Header: header}
	return tsav.Cookies()
}

// Network request builder.
func Network(url, method, query string) (req *http.Request, err error) {
	switch method {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
		req.URL.RawQuery = query
	case "POST":
		req, err = http.NewRequest("POST", url, strings.NewReader(query))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	return
}

// Proc json function in order to reduce code.
func JsonProc(body *http.Response, container interface{}) error {
	if err := json.NewDecoder(body.Body).Decode(container); err != nil {
		return err
	}
	return nil
}

// Build up fake video page referer
func BuildVideoReferer(id int) string {
	return fmt.Sprintf(VIDEO_URL, id)
}
