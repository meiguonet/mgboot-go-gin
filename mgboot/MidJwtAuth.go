package mgboot

import (
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"github.com/meiguonet/mgboot-go-common/enum/RegexConst"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"github.com/meiguonet/mgboot-go-fiber/enum/JwtVerifyErrno"
	"strings"
)

func MidJwtAuth(settingsKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidJwtAuth")
		}

		if settingsKey == "" {
			ctx.Next()
			return
		}

		settings := GetJwtSettings(settingsKey)

		if settings == nil {
			ctx.Next()
			return
		}

		token := strings.TrimSpace(ctx.GetHeader("Authorization"))
		token = stringx.RegexReplace(token, RegexConst.SpaceSep, " ")

		if strings.Contains(token, " ") {
			token = stringx.SubstringAfter(token, " ")
		}

		if token == "" {
			panic(NewJwtAuthError(JwtVerifyErrno.NotFound))
		}

		errno := VerifyJsonWebToken(token, settings)

		if errno < 0 {
			panic(NewJwtAuthError(errno))
		}

		ctx.Next()
	}
}
