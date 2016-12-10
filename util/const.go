package util

import (
	"log"
	"net/url"
)

const (
	USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36"
)

var (
	LOGIN_ERR_MAP = map[int]string{
		-105: "验证码错误",
		-618: "昵称重复或含有非法字符",
		-619: "昵称不能小于3个字符或者大于30个字符",
		-620: "该昵称已被使用",
		-622: "Email已存在",
		-625: "密码错误次数过多",
		-626: "用户不存在",
		-627: "密码错误",
		-628: "密码不能小于6个字符或大于16个字符",
		-645: "昵称或密码过短",
		-646: "请输入正确的手机号",
		-647: "该手机已绑定另外一个账号",
		-648: "验证码发送失败",
		-652: "历史遗留问题，昵称与手机号重复，请联系管理员",
		-662: "加密后的密码已过期",
	}
)

var BiliLoginURL *url.URL

func init() {
	var err error
	BiliLoginURL, err = url.Parse("https://passport.bilibili.com/ajax/miniLogin/login")
	if err != nil {
		log.Fatal(err)
	}
}
