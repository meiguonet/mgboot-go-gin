package mgboot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/numberx"
	"github.com/meiguonet/mgboot-go-common/util/slicex"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"strings"
	"time"
)

var runtimeLogger logx.Logger
var requestLogLogger logx.Logger
var logRequestBody bool
var executeTimeLogLogger logx.Logger
var errorHandlers = make([]ErrorHandler, 0)

func RuntimeLogger(logger ...logx.Logger) logx.Logger {
	if len(logger) > 0 {
		runtimeLogger = logger[0]
	}

	l := runtimeLogger

	if l == nil {
		l = NewNoopLogger()
	}

	return l
}

func RequestLogLogger(logger ...logx.Logger) logx.Logger {
	if len(logger) > 0 {
		requestLogLogger = logger[0]
	}

	l := requestLogLogger

	if l == nil {
		l = NewNoopLogger()
	}

	return l
}

func RequestLogEnabled() bool {
	return requestLogLogger != nil
}

func LogRequestBody(flag ...bool) bool {
	if len(flag) > 0 {
		logRequestBody = flag[0]
	}

	return logRequestBody
}

func ExecuteTimeLogLogger(logger ...logx.Logger) logx.Logger {
	if len(logger) > 0 {
		executeTimeLogLogger = logger[0]
	}

	l := executeTimeLogLogger

	if l == nil {
		l = NewNoopLogger()
	}

	return l
}

func ExecuteTimeLogEnabled() bool {
	return executeTimeLogLogger != nil
}

func LogExecuteTime(ctx *gin.Context) {
	if !ExecuteTimeLogEnabled() {
		return
	}

	req := NewRequest(ctx)
	elapsedTime := calcElapsedTime(ctx)

	if elapsedTime == "" {
		return
	}

	sb := strings.Builder{}
	sb.WriteString(req.GetMethod())
	sb.WriteString(" ")
	sb.WriteString(req.GetRequestUrl(true))
	sb.WriteString(", total elapsed time: " + elapsedTime)
	ExecuteTimeLogLogger().Info(sb.String())
	ctx.Set("X-Response-Time", elapsedTime)
}

func WithBuiltinErrorHandlers() {
	errorHandlers = []ErrorHandler{
		NewRateLimitErrorHandler(),
		NewJwtAuthErrorHandler(),
		NewValidateErrorHandler(),
	}
}

func ReplaceBuiltinErrorHandler(errName string, handler ErrorHandler) {
	errName = stringx.EnsureRight(errName, "Error")
	errName = stringx.EnsureLeft(errName, "builtin.")
	handlers := make([]ErrorHandler, 0)
	var added bool

	for _, h := range errorHandlers {
		if h.GetErrorName() == errName {
			handlers = append(handlers, handler)
			added = true
			continue
		}

		handlers = append(handlers, h)
	}

	if !added {
		handlers = append(handlers, handler)
	}

	errorHandlers = handlers
}

func WithErrorHandler(handler ErrorHandler) {
	handlers := make([]ErrorHandler, 0)
	var added bool

	for _, h := range errorHandlers {
		if h.GetErrorName() == handler.GetErrorName() {
			handlers = append(handlers, handler)
			added = true
			continue
		}

		handlers = append(handlers, h)
	}

	if !added {
		handlers = append(handlers, handler)
	}

	errorHandlers = handlers
}

func WithErrorHandlers(handlers []ErrorHandler) {
	if len(handlers) < 1 {
		return
	}

	for _, handler := range handlers {
		WithErrorHandler(handler)
	}
}

func ErrorHandlers() []ErrorHandler {
	return errorHandlers
}

func NeedCorsSupport(ctx *gin.Context) bool {
	req := NewRequest(ctx)
	methods := []string{"PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	if slicex.InStringSlice(req.GetMethod(), methods) {
		return true
	}

	contentType := strings.ToLower(req.GetHeader("Content-Type"))

	if strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data") ||
		strings.Contains(contentType, "text/plain") {
		return true
	}

	headerNames := []string{
		"Accept",
		"Accept-Language",
		"Content-Language",
		"DPR",
		"Downlink",
		"Save-Data",
		"Viewport-Widt",
		"Width",
	}

	for headerName := range req.GetHeaders() {
		if slicex.InStringSlice(headerName, headerNames) {
			return true
		}
	}

	return false
}

func AddCorsSupport(ctx *gin.Context) {
	if !NeedCorsSupport(ctx) {
		return
	}
	
	settings := GetCorsSettings()
	
	if settings == nil {
		return
	}

	allowedOrigins := settings.AllowedOrigins()

	if slicex.InStringSlice("*", allowedOrigins) {
		ctx.Header("Access-Control-Allow-Origin", "*")
	} else {
		ctx.Header("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))
	}

	allowedHeaders := settings.AllowedHeaders()

	if len(allowedHeaders) > 0 {
		ctx.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	}

	exposedHeaders := settings.ExposedHeaders()

	if len(exposedHeaders) > 0 {
		ctx.Header("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
	}

	maxAge := settings.MaxAge()

	if maxAge > 0 {
		n1 := castx.ToInt64(maxAge.Seconds())
		ctx.Header("Access-Control-Max-Age", fmt.Sprintf("%d", n1))
	}

	if settings.AllowCredentials() {
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}
}

func AddPoweredBy(ctx *gin.Context) {
	poweredBy := AppConf.GetString("app.poweredBy")

	if poweredBy == "" {
		return
	}

	ctx.Header("X-Powered-By", poweredBy)
}

func calcElapsedTime(ctx *gin.Context) string {
	var execStart time.Time
	v1, _ := ctx.Get("ExecStart")

	if d1, ok := v1.(time.Time); ok {
		execStart = d1
	}

	if execStart.IsZero() {
		return ""
	}

	d2 := time.Now().Sub(execStart)

	if d2 < time.Second {
		return fmt.Sprintf("%dms", d2)
	}

	n1 := d2.Seconds()
	return numberx.ToDecimalString(n1, 3) + "ms"
}
