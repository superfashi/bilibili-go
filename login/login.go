package bililogin

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/hanbang-wang/bilibili-go/util"
	"github.com/skratchdot/open-golang/open"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type rsaLogin struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

type userAccess struct {
	Status  bool `json:"status"`
	Message struct {
		Code int `json:"code"`
	} `json:"message"`
}

// Log into bilibili with credential.
func Login(client *http.Client, username, password string) error {
	// Fake requests to main & login page
	if req, err := util.Network("http://www.bilibili.com", "GET", ""); err == nil {
		if _, err = client.Do(req); err != nil {
			return err
		}
	} else {
		return err
	}
	if req, err := util.Network("https://passport.bilibili.com/login", "GET", ""); err == nil {
		if _, err = client.Do(req); err != nil {
			return err
		}
	} else {
		return err
	}
	// Get captcha
	vdcode, err := getCaptcha(client)
	if err != nil {
		return err
	}

	// Get encoded password
	pass, err := rsaEncrypt(client, []byte(password))
	if err != nil {
		return err
	}

	// Build login parameters
	data := url.Values{}
	data.Add("captcha", vdcode)
	data.Add("userid", username)
	data.Add("pwd", pass)
	data.Add("keep", "1")
	req, err := util.Network(util.BiliLoginURL.String(), "POST", data.Encode())
	if err != nil {
		return err
	}
	req.Header.Set("Origin", "https://passport.bilibili.com")
	req.Header.Set("Referer", "https://passport.bilibili.com/ajax/miniLogin/minilogin")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	user := new(userAccess)
	if err = util.JsonProc(resp, user); err != nil {
		return err
	}
	if user.Status {
		return nil
	} else {
		if info, ok := util.LOGIN_ERR_MAP[user.Message.Code]; ok {
			return errors.New(info)
		} else {
			return fmt.Errorf("Unknown error with code: %d", user.Message.Code)
		}
	}
}

func getCaptcha(client *http.Client) (string, error) {
	var ret string
	req, err := util.Network("https://passport.bilibili.com/captcha", "GET", "")
	if err != nil {
		return "", err
	}
	req.Header.Set("Referer", "https://passport.bilibili.com/ajax/miniLogin/minilogin")
	req.Header.Set("Accept", "image/jpeg")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	tmpjpg := filepath.Join(os.TempDir(), "vdcode.jpg")
	tmpfil, err := os.Create(tmpjpg)
	if err != nil {
		return "", err
	}
	defer syscall.Unlink(tmpjpg)
	defer tmpfil.Close()

	if _, err = io.Copy(tmpfil, resp.Body); err != nil {
		return "", err
	}
	tmpfil.Close()

	if err = open.Start(tmpjpg); err != nil {
		return "", err
	}

	fmt.Print("请输入你看到的验证码并回车：")
	fmt.Scanf("%s", &ret)
	return strings.ToLower(ret), nil
}

func rsaEncrypt(client *http.Client, data []byte) (string, error) {
	ret := new(rsaLogin)
	req, err := util.Network("https://passport.bilibili.com/login", "GET", "act=getkey")
	if err != nil {
		return "", err
	}
	req.Header.Set("Referer", "https://passport.bilibili.com/ajax/miniLogin/minilogin")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if err = util.JsonProc(resp, ret); err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(ret.Key))
	if block == nil {
		return "", errors.New("Error reading public key.")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	encrypted, err := rsa.EncryptPKCS1v15(crand.Reader, pubInterface.(*rsa.PublicKey), append([]byte(ret.Hash), data...))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encrypted), nil
}
