package mgboot

import (
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	"github.com/meiguonet/mgboot-go-common/util/validatex"
	"strings"
)

func MidValidate(arg0 interface{}) gin.HandlerFunc {
	rules := make([]string, 0)
	var failfast bool

	if items, ok := arg0.([]string); ok && len(items) > 0 {
		for _, s1 := range items {
			if s1 == "" || s1 == "false" {
				continue
			}

			if s1 == "true" {
				failfast = true
				continue
			}

			rules = append(rules, s1)
		}
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		s1 = strings.ReplaceAll(s1, "[syh]", `"`)
		entries := jsonx.ArrayFrom(s1)

		for _, entry := range entries {
			s2, ok := entry.(string)

			if !ok || s2 == "" || s2 == "false" {
				continue
			}

			if s2 == "true" {
				failfast = true
				continue
			}

			rules = append(rules, s2)
		}
	}

	return func(ctx *gin.Context) {
		if AppConf.GetBoolean("logging.logMiddlewareRun") {
			RuntimeLogger().Info("middleware run: mgboot.MidValidate")
		}

		if len(rules) < 1 {
			ctx.Next()
			return
		}

		validator := validatex.NewValidator()
		req := NewRequest(ctx)
		data := req.GetMap()

		if failfast {
			errorTips := validatex.FailfastValidate(validator, data, rules)

			if errorTips != "" {
				panic(NewValidateError(errorTips, true))
			}

			ctx.Next()
			return
		}

		validateErrors := validatex.Validate(validator, data, rules)

		if len(validateErrors) > 0 {
			panic(NewValidateError(validateErrors))
		}

		ctx.Next()
	}
}
