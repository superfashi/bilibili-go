package util

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func errorP(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func StringToCookie(cookie string) []*http.Cookie {
	tsav := http.Request{}
	tsav.Header.Add("Cookie", cookie)
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
