package biligo

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hanbang-wang/bilibili-go/util"
)

// IsLoggedIn determines whether you logged in
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

// SendComment sends a comment, returns success or not.
// Requesting too quick will cause captcha requirement.
func (b *Bilibili) SendComment(oid int, message string) (int, error) {
	ret := new(util.AddComment)

	// Build request parameters
	data := url.Values{}
	data.Add("jsonp", "jsonp")
	data.Add("message", message)
	data.Add("type", "1")
	data.Add("plat", "1")
	data.Add("oid", strconv.Itoa(oid))

	req, err := util.Network(util.COMMENT_URL, "POST", data.Encode())
	if err != nil {
		return -1, err
	}
	req.Header.Add("Origin", util.MAIN_HOST)
	req.Header.Add("Referer", util.BuildVideoReferer(oid))
	resp, err := b.Client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	if err = util.JsonProc(resp, ret); err != nil {
		return -1, err
	}
	if ret.Code != 0 {
		return ret.Code, fmt.Errorf("Unknown error! Server returned %d with message %s.", ret.Code, ret.Message)
	}
	return 0, nil
}
