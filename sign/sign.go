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
	if o.Sign == "" {
		panic(errors.New("sign is required"))
	}
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
	accessKey := ctx.GetHeader("access_key")
	if accessKey != SharedOptions.AccessKey {
		log.Errorf("Error accesskey %s\n", accessKey)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedError())
		return
	}
	ctx.Next()
}
