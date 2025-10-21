package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ydxlt/go-foundation/dto"
	"github.com/ydxlt/go-foundation/errs"
	"github.com/ydxlt/go-foundation/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Claims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

// Options 配置结构
type Options struct {
	JwtKey    string        `mapstructure:"jwt_key"`
	ExpiresAt time.Duration `mapstructure:"expires_at"`
	Issuer    string        `mapstructure:"issuer"`
}

func (o *Options) Validate() {
	// 手动校验必填字段
	if o.JwtKey == "" {
		panic("token.jwt_key is required")
	}
	if o.ExpiresAt == 0 {
		panic("token.expires_at is required")
	}
	if o.Issuer == "" {
		panic("token.issuer is required")
	}
}

var OpenWhitePaths []string

var (
	cfg  *Options
	once sync.Once
)

func Init(c *Options) {
	once.Do(func() {
		cfg = c
		log.Infof("[token] InitWithConfig loaded: %+v", cfg)
	})
}

func loadFromViper() {
	once.Do(func() {
		var c Options
		if err := viper.Sub("token").Unmarshal(&c); err != nil {
			log.Infof("[token] viper unmarshal error: %v", err)
			panic(err)
		}
		c.Validate()
		cfg = &c
		log.Infof("[token] InitFromViper loaded: %+v", cfg)
	})
}

func configs() Options {
	if cfg == nil {
		loadFromViper() // 默认走 viper
	}
	return *cfg
}

// CheckToken 校验中间件
func CheckToken(ctx *gin.Context) {
	for _, p := range OpenWhitePaths {
		if ctx.Request.URL.Path == p {
			ctx.Next()
			return
		}
	}
	tokenString := ctx.GetHeader(dto.HeaderAuthorization)
	uid, err := ValidateAccessToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedError())
		return
	}
	ctx.Set(dto.HeaderUID, uid)
	log.Debugf("CheckToken success, uid = %d", uid)
	ctx.Next()
}

// GenerateAccessToken 生成新的 Access AccessToken
func GenerateAccessToken(uid int64) (string, int64, error) {
	now := time.Now()
	expiration := now.Add(configs().ExpiresAt)

	claims := Claims{
		ID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "bk",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(configs().JwtKey))
	if err != nil {
		return "", 0, err
	}
	return token, expiration.Unix(), nil
}

// GenerateRefreshToken 生成新的 Refresh AccessToken（有效期30天）
func GenerateRefreshToken() (string, time.Time, error) {
	refreshToken, err := generateSecureRandomString(32)
	if err != nil || refreshToken == "" {
		log.Debugf("ErrorWithInternal generating refresh token: %v", err)
		return "", time.Time{}, err
	}
	expireAt := time.Now().Add(24 * 30 * time.Hour) // 直接返回 time.Time
	return refreshToken, expireAt, nil
}

// generateSecureRandomString 生成安全的随机字符串
func generateSecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func GenerateTokens(uid int64) (Tokens, error) {
	accessToken, _, err := GenerateAccessToken(uid)
	if err != nil {
		log.Errorf("GenerateTokens err %s", err)
		return Tokens{}, errs.InternalError
	}
	refreshToken, rTokenExpiredAt, err := GenerateRefreshToken()
	if err != nil {
		log.Errorf("GenerateTokens err %s", err)
		return Tokens{}, err
	}
	return Tokens{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RTokenExpireAt: rTokenExpiredAt,
	}, nil
}

func ValidateAccessToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return configs().JwtKey, nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, fmt.Errorf("token expired: %w", err) // 明确返回过期错误
			}
		}
		return 0, fmt.Errorf("token invalid: %w", err) // 其他错误
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, errors.New("token claims invalid")
	}

	return claims.ID, nil
}

func GetIDFromToken(accessToken string) (int64, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &Claims{})
	if err != nil {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.ID, nil
}

func GetUIDFromContext(ctx *gin.Context) (int64, error) {
	uid := ctx.GetInt64(dto.HeaderUID)
	if uid > 0 {
		return uid, nil
	}
	var err error
	tokenString := ctx.GetHeader("Authorization")
	uid, err = GetIDFromToken(tokenString)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

type Tokens struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	RTokenExpireAt time.Time `json:"-"`
}
