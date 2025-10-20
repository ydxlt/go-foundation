package dto

import (
	"github.com/gin-gonic/gin"
)

var (
	HeaderAuthorization = "Authorization"
	HeaderAppID         = "X-App-ID"
	HeaderAppSign       = "X-App-Sign"
	HeaderAccessKey     = "X-Access-Key"
	HeaderUID           = "X-User-ID"
	HeaderCID           = "X-Client-ID"
	HeaderVC            = "X-Version-Code"
	HeaderVN            = "X-Version-Name"
	HeaderOS            = "X-Device-OS"
	HeaderModel         = "X-Device-Model"
	HeaderBrand         = "X-Device-Brand"
	HeaderLanguage      = "X-Device-Language"
)

type CommHeaders struct {
	VC       int64  `json:"X-Version-Code"`
	VN       string `json:"X-Version-Name"`
	AppID    string `json:"X-App-ID"`
	Sign     string `json:"X-App-Sign"`
	UID      int64  `json:"X-User-ID"`
	ClientID string `json:"X-Client-ID"`
	OS       string `json:"X-Device-OS"`
	Model    string `json:"X-Device-Model"`
	Brand    string `json:"X-Device-Brand"`
	Language string `json:"X-Device-Language"`
}

type legacyCommHeaders struct {
	VC       int64  `json:"vc"`
	VN       string `json:"vn"`
	Package  string `json:"package" binding:"required"`
	Cid      string `json:"cid"`
	Aid      string `json:"aid"`
	Sign     string `json:"sign" binding:"required"`
	Uid      int64  `json:"uid"`
	OS       string `json:"os"`
	Model    string `json:"model"`
	Brand    string `json:"brand"`
	Language string `json:"language"`
}

func BindHeader(ctx *gin.Context) (CommHeaders, error) {
	// 先尝试绑定LegacyHeader，如果失败再绑定新的，统一返回新的Header
	var legacyHeader legacyCommHeaders
	if err := ctx.ShouldBindHeader(&legacyHeader); err == nil {
		return CommHeaders{
			VC:       legacyHeader.VC,
			VN:       legacyHeader.VN,
			AppID:    legacyHeader.Aid,
			Sign:     legacyHeader.Sign,
			ClientID: legacyHeader.Cid,
			UID:      legacyHeader.Uid,
			OS:       legacyHeader.OS,
			Model:    legacyHeader.Model,
			Brand:    legacyHeader.Brand,
			Language: legacyHeader.Language,
		}, nil
	}

	var header CommHeaders
	if err := ctx.ShouldBindHeader(&header); err != nil {
		return CommHeaders{}, err
	}

	return header, nil
}
