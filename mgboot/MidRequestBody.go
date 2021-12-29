package mgboot

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"time"
)

func MidRequestBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("ExecStart", time.Now())
		req := NewRequest(ctx)
		method := req.GetMethod()

		if method == "GET" {
			ctx.Next()
			return
		}

		contentType := req.GetHeader("Content-Type")
		isJsonPayload := strings.Contains(contentType, "json")
		isXmlPayload := strings.Contains(contentType, "xml")

		if !isJsonPayload && !isXmlPayload {
			ctx.Next()
			return
		}

		var buf []byte

		if ctx.Request.Body != nil {
			buf, _ = ctx.GetRawData()
		}

		if len(buf) < 1 {
			buf = []byte{}
		}

		ctx.Set("requestRawBody", buf)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(buf))
		ctx.Next()
	}
}
