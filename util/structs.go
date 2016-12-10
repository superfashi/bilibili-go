package util

type UserInfo struct {
	Code   int  `json:"code"`
	Status bool `json:"status"`
	Data   struct {
		LevelInfo struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		BCoins           int     `json:"bCoins"`
		Coins            float64 `json:"coins"`
		Face             string  `json:"face"`
		NameplateCurrent string  `json:"nameplate_current"`
		Uname            string  `json:"uname"`
		UserStatus       string  `json:"userStatus"`
		VipType          int     `json:"vipType"`
		VipStatus        int     `json:"vipStatus"`
		OfficialVerify   int     `json:"official_verify"`
	} `json:"data"`
}

type AddComment struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
