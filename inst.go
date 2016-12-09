package biligo

import (
	"github.com/hanbang-wang/bilibili-go/login"
	"github.com/hanbang-wang/bilibili-go/util"
	"net/http"
	"net/http/cookiejar"
)

type Bilibili struct {
	Client *http.Client
}

func NewFromLogin(username, password string) (*Bilibili, error) {
	ret := &Bilibili{}
	err := bililogin.LoginWith(ret.Client, username, password)
	return ret, err
}

func NewFromCookie(cookie string) (*Bilibili, error) {
	ret := &Bilibili{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(util.BilibiliURL, util.StringToCookie(cookie))
	ret.Client = &http.Client{Jar: jar}
	return ret, nil
}
