package mgboot

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/meiguonet/mgboot-go-common/AppConf"
)

func MidFinalStep() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidFinalResult")
		}

		LogExecuteTime(ctx)
		AddCorsSupport(ctx)
		AddPoweredBy(ctx)
		v1, _ := ctx.Get("ResponsePayload")
		var payload ResponsePayload

		if p, ok := v1.(ResponsePayload); ok {
			payload = p
		}

		if payload == nil {
			ctx.Render(200, render.Data{
				ContentType: "text/html; charset=utf-8",
				Data:        []byte("unsupported response payload found"),
			})

			return
		}

		statusCode, contents := payload.GetContents()

		if statusCode >= 400 {
			ctx.AbortWithStatus(statusCode)
			return
		}

		if pl, ok := payload.(AttachmentResponse); ok {
			pl.AddSpecifyHeaders(ctx)

			ctx.Render(200, render.Data{
				ContentType: pl.GetContentType(),
				Data:        pl.Buffer(),
			})

			return
		}

		if pl, ok := payload.(ImageResponse); ok {
			ctx.Render(200, render.Data{
				ContentType: pl.GetContentType(),
				Data:        pl.Buffer(),
			})

			return
		}

		ctx.Render(200, render.Data{
			ContentType: payload.GetContentType(),
			Data:        []byte(contents),
		})
	}
}
