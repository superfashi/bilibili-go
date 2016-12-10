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

// Login Bilibili with credentials.
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

// Login Bilibili with cookie
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

// Determine whether you logged in
func (b *Bilibili) IsLoggedIn() (bool, error) {
	req, err := util.Network("https://account.bilibili.com/home/userInfo", "GET", "")
	if err != nil {
		return false, err
	}
	resp, err := b.Client.Do(req)
	userinfo := new(util.UserInfo)
	if err = util.JsonProc(resp, userinfo); err != nil {
		return false, err
	}
	return userinfo.Code == 0, nil
}
