package mgboot

import (
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"strings"
)

func MidRequestLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidRequestLog")
		}

		if !RequestLogEnabled() {
			ctx.Next()
			return
		}

		req := NewRequest(ctx)
		logger := RequestLogLogger()
		sb := strings.Builder{}
		sb.WriteString(req.GetMethod())
		sb.WriteString(" ")
		sb.WriteString(req.GetRequestUrl(true))
		sb.WriteString(" from ")
		sb.WriteString(req.GetClientIp())
		logger.Info(sb.String())

		if LogRequestBody() {
			rawBody := req.GetRawBody()

			if len(rawBody) > 0 {
				logger.Debugf(string(rawBody))
			}
		}

		ctx.Next()
	}
}
