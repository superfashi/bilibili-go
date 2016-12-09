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
	"syscall"
)

type rsaLogin struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

type userAccess struct {
	Status bool `json:"status"`
	Data   struct {
		Code int `json:"code"`
	} `json:"data"`
}

// Log into bilibili with credential.
func LoginWith(client *http.Client, username, password string) error {
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
	if req, err := util.Network("https://passport.bilibili.com/ajax/miniLogin/login", "POST", data.Encode()); err == nil {
		req.Header.Set("Origin", "https://passport.bilibili.com")
		req.Header.Set("Referer", "https://passport.bilibili.com/ajax/miniLogin/minilogin")
		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()
			user := new(userAccess)
			if err = util.JsonProc(resp, user); err == nil {
				if user.Status {
					return nil
				} else {
					if info, ok := util.LOGIN_ERR_MAP[user.Data.Code]; ok {
						return errors.New(info)
					} else {
						return fmt.Errorf("Unknown error with code: %d", user.Data.Code)
					}
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func getCaptcha(client *http.Client) (string, error) {
	var ret string
	if req, err := util.Network("https://passport.bilibili.com/captcha", "GET", ""); err == nil {
		if resp, err := client.Do(req); err == nil {
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

			return ret, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
	return "", nil // Impossible
}

func rsaEncrypt(client *http.Client, data []byte) (string, error) {
	ret := new(rsaLogin)
	if req, err := util.Network("https://passport.bilibili.com/login", "GET", "act=getkey"); err == nil {
		if resp, err := client.Do(req); err == nil {
			if err = util.JsonProc(resp, ret); err == nil {
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
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
	return "", nil // Impossible
}
