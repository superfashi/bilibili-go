package biligo

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/hanbang-wang/bilibili-go/login"
	"github.com/hanbang-wang/bilibili-go/util"
)

// Bilibili is a struct for easy net.Client access
type Bilibili struct {
	Client *http.Client
}

// NewFromLogin logs into Bilibili with credentials.
func NewFromLogin(username, password string) (*Bilibili, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	ret := &Bilibili{&http.Client{Jar: jar}}
	if err = bililogin.Login(ret.Client, username, password); err != nil {
		return nil, err
	}
	return ret, nil
}

// NewFromCookie log into Bilibili with cookie
// Strongly deprecated
func NewFromCookie(cookie string) (*Bilibili, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	procCookie := util.StringToCookie(cookie)
	for i := range procCookie {
		procCookie[i].Domain = ".bilibili.com"
		procCookie[i].Path = "/"
	}
	jar.SetCookies(util.BiliLoginURL, procCookie)
	ret := &Bilibili{&http.Client{Jar: jar}}
	return ret, nil
}
