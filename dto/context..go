package dto

import (
	"go-api-comm/errs"
	"go-api-comm/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidateAndGetUID 公共校验方法，返回uid和是否通过校验
func ValidateAndGetUID(ctx *gin.Context) (int64, bool) {
	uid, err := GetUIDFromContext(ctx)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedError())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, InternalError())
		}
		return 0, false
	}
	return uid, true
}

func GetUIDFromContext(ctx *gin.Context) (int64, error) {
	uidStr := ctx.Request.Header.Get("uid")
	if uidStr == "" {
		log.Warnf("uid parse error, uid: %v", uidStr)
		// 处理没有 UID 的情况
		return 0, errs.UnauthorizedError
	}

	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid == 0 {
		log.Warnf("uid parse error, uid: %v", uidStr)
		return 0, errs.InvalidParamError
	}
	return uid, nil
}
