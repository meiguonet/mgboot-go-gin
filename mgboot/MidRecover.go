package mgboot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"github.com/meiguonet/mgboot-go-common/util/errorx"
)

func MidRecover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidRecover")
		}

		defer func() {
			r := recover()

			if r == nil {
				return
			}

			var err error

			if ex, ok := r.(error); ok {
				err = ex
			} else {
				err = fmt.Errorf("%v", r)
			}

			if err == nil {
				return
			}

			handlers := ErrorHandlers()
			var handler ErrorHandler

			for _, h := range handlers {
				if h.MatchError(err) {
					handler = h
					break
				}
			}

			LogExecuteTime(ctx)
			AddPoweredBy(ctx)
			AddCorsSupport(ctx)

			if handler == nil {
				RuntimeLogger().Error(errorx.Stacktrace(err))
				ctx.AbortWithStatus(500)
				return
			}

			if ex, ok := err.(RateLimitError); ok {
				ex.AddSpecifyHeaders(ctx)
			}

			payload := handler.HandleError(err)
			statusCode, contents := payload.GetContents()

			if statusCode >= 400 {
				ctx.AbortWithStatus(statusCode)
				return
			}

			ctx.Render(200, render.Data{
				ContentType: payload.GetContentType(),
				Data:        []byte(contents),
			})
		}()

		ctx.Next()
	}
}
