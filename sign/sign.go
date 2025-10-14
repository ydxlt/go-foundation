package sign

import (
	"net/http"

	"github.com/ydxlt/go-foundation/dto"
	"github.com/ydxlt/go-foundation/log"

	"github.com/gin-gonic/gin"
)

type Configs struct {
	Sign      string
	AccessKey string
	WhiteList []string
}

var configs *Configs

func requiredConfigs() {
	if configs == nil {
		panic("configs is nil")
	}
}

func CheckSign(ctx *gin.Context) {
	requiredConfigs()
	for _, p := range (*configs).WhiteList {
		if ctx.Request.URL.Path == p {
			ctx.Next()
			return
		}
	}
	sign := ctx.GetHeader("sign")
	if sign != configs.Sign {
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
	if accessKey != configs.AccessKey {
		log.Errorf("Error accesskey %s\n", accessKey)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedError())
		return
	}
	ctx.Next()
}
