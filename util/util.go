package util

import (
	"encoding/json"
	"net/http"
	"strings"
)

func StringToCookie(cookie string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", cookie)
	tsav := &http.Request{Header: header}
	return tsav.Cookies()
}

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

func JsonProc(body *http.Response, container interface{}) error {
	if err := json.NewDecoder(body.Body).Decode(container); err != nil {
		return err
	}
	return nil
}
