package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ydxlt/go-foundation/log"
)

var (
	HeaderAuthorization = "Authorization"
	HeaderAppID         = "X-App-ID"
	HeaderAppSign       = "X-App-Sign"
	HeaderAppChannel    = "X-App-Channel"
	HeaderAccessKey     = "X-Access-Key"
	HeaderUID           = "X-User-ID"
	HeaderCID           = "X-Client-ID"
	HeaderVC            = "X-Version-Code"
	HeaderVN            = "X-Version-Name"
	HeaderOS            = "X-Device-OS"
	HeaderOSVer         = "X-Device-OS-Version"
	HeaderModel         = "X-Device-Model"
	HeaderBrand         = "X-Device-Brand"
	HeaderLanguage      = "X-Device-Language"
	HeaderPlatform      = "X-Device-Platform"
)

type CommHeaders struct {
	VC       int64  `header:"X-Version-Code"`
	VN       string `header:"X-Version-Name"`
	AppID    string `header:"X-App-ID"`
	Sign     string `header:"X-App-Sign"`
	UID      int64  `header:"X-User-ID"`
	ClientID string `header:"X-Client-ID"`
	OS       string `header:"X-Device-OS"`
	OSVer    string `header:"X-Device-OS-Version"`
	Model    string `header:"X-Device-Model"`
	Brand    string `header:"X-Device-Brand"`
	Language string `header:"X-Device-Language"`
}

func (h *CommHeaders) DeviceInfo() string {
	data := map[string]string{
		"OS":       h.OS,
		"OSVer":    h.OSVer,
		"Model":    h.Model,
		"Brand":    h.Brand,
		"Language": h.Language,
	}
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

type legacyCommHeaders struct {
	VC       int64  `header:"vc"`
	VN       string `header:"vn"`
	Package  string `header:"package" binding:"required"`
	Cid      string `header:"cid"`
	Aid      string `header:"aid"`
	Sign     string `header:"sign" binding:"required"`
	Uid      int64  `header:"uid"`
	OS       string `header:"os"`
	OSVer    string `header:"osVerson"`
	Model    string `header:"model"`
	Brand    string `header:"brand"`
	Language string `header:"language"`
}

func BindHeader(ctx *gin.Context) (CommHeaders, error) {
	// 先尝试绑定LegacyHeader，如果失败再绑定新的，统一返回新的Header
	var legacyHeader legacyCommHeaders
	if err := ctx.ShouldBindHeader(&legacyHeader); err == nil {
		log.Debugf("bind legacyHeader: %v", legacyHeader)
		return CommHeaders{
			VC:       legacyHeader.VC,
			VN:       legacyHeader.VN,
			AppID:    legacyHeader.Aid,
			Sign:     legacyHeader.Sign,
			ClientID: legacyHeader.Cid,
			UID:      legacyHeader.Uid,
			OS:       legacyHeader.OS,
			OSVer:    legacyHeader.OSVer,
			Model:    legacyHeader.Model,
			Brand:    legacyHeader.Brand,
			Language: legacyHeader.Language,
		}, nil
	}

	var header CommHeaders
	if err := ctx.ShouldBindHeader(&header); err != nil {
		log.Errorf("bind header err: %v", err)
		return CommHeaders{}, err
	}
	log.Debugf("bind header success %+v", header)
	return header, nil
}
