package sign

import (
	"errors"
	"net/http"

	"github.com/ydxlt/go-foundation/dto"
	"github.com/ydxlt/go-foundation/log"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Sign      string
	AccessKey string
	WhiteList []string
}

func (o *Options) validate() {
	if o.AccessKey == "" {
		panic(errors.New("access_key is required"))
	}
}

var SharedOptions *Options

func requiredConfigs() {
	if SharedOptions == nil {
		panic("configs is nil")
	}
	SharedOptions.validate()
}

func CheckSign(ctx *gin.Context) {
	requiredConfigs()
	for _, p := range (*SharedOptions).WhiteList {
		if ctx.Request.URL.Path == p {
			ctx.Next()
			return
		}
	}
	sign := ctx.GetHeader(dto.HeaderAppSign)
	if sign == "" {
		sign = ctx.GetHeader("sign") // 兼容老版本
	}
	if sign != SharedOptions.Sign {
		log.Errorf("Error sign： %s\n", sign)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedError())
		return
	}
	ctx.Next()
}

func CheckAccessKey(ctx *gin.Context) {
	requiredConfigs()
	// access_key获取不到?
	accessKey := ctx.GetHeader(dto.HeaderAccessKey)
	if accessKey == "" {
		accessKey = ctx.GetHeader("access_key") // 兼容老版本
	}
	if accessKey != SharedOptions.AccessKey {
		log.Errorf("Error %s %s\n", dto.HeaderAccessKey, accessKey)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedError())
		return
	}
	ctx.Next()
}
