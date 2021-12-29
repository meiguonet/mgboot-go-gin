package mgboot

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"strings"
)

func MidOptionsReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidOptionsReq")
		}

		if strings.ToUpper(ctx.Request.Method) != "OPTIONS" {
			ctx.Next()
			return
		}

		AddCorsSupport(ctx)
		AddPoweredBy(ctx)

		ctx.Render(200, render.Data{
			ContentType: "application/json; charset=utf-8",
			Data:        []byte(`{"code":200}`),
		})
	}
}
