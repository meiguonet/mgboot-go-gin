package mgboot

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/enum/RegexConst"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	"github.com/meiguonet/mgboot-go-common/util/mapx"
	"github.com/meiguonet/mgboot-go-common/util/slicex"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"math"
	"mime/multipart"
	"net/url"
	"regexp"
	"strings"
)

type Request struct {
	ctx *gin.Context
}

func NewRequest(ctx *gin.Context) *Request {
	return &Request{ctx: ctx}
}

func (r *Request) GetMethod() string {
	return strings.ToUpper(r.ctx.Request.Method)
}

func (r *Request) GetHeaders() map[string]string {
	if len(r.ctx.Request.Header) < 1 {
		return map[string]string{}
	}

	headers := map[string]string{}

	for name, values := range r.ctx.Request.Header {
		if len(values) < 1 {
			headers[name] = ""
			continue
		}

		headers[name] = strings.Join(values, ",")
	}

	return headers
}

func (r *Request) GetHeader(name string) string {
	name = strings.ToLower(name)
	headers := r.GetHeaders()

	for headerName, headerValue := range headers {
		if strings.ToLower(headerName) == name {
			return headerValue
		}
	}

	return ""
}

func (r *Request) GetQueryParams() map[string]string {
	map1 := map[string]string{}
	query := r.ctx.Request.URL.RawQuery

	if query == "" {
		return map1
	}

	queryMap, err := url.ParseQuery(query)

	if err != nil {
		return map1
	}

	for key, values := range queryMap {
		if key == "" || len(values) < 1 {
			continue
		}

		map1[key] = values[0]
	}

	return map1
}

func (r *Request) GetQueryString(urlencode ...bool) string {
	params := r.GetQueryParams()

	if len(params) < 1 {
		return ""
	}

	if len(urlencode) > 0 && urlencode[0] {
		values := url.Values{}

		for name, value := range params {
			values[name] = []string{value}
		}

		return values.Encode()
	}

	sb := strings.Builder{}
	n1 := 0

	for name, value := range params {
		sb.WriteString(name + "=" + value)

		if n1 > 0 {
			sb.WriteString("&")
		}

		n1++
	}

	return sb.String()
}

func (r *Request) GetRequestUrl(withQueryString ...bool) string {
	s1 := r.ctx.Request.URL.RequestURI()
	s1 = stringx.EnsureLeft(s1, "/")

	if len(withQueryString) > 0 && withQueryString[0] {
		qs := r.GetQueryString()

		if qs != "" {
			s1 += "?" + qs
		}
	}

	return s1
}

func (r *Request) GetFormData() map[string]string {
	map1 := map[string]string{}
	r.ctx.PostForm("NonExistsKey")

	if len(r.ctx.Request.PostForm) < 1 {
		return map1
	}

	for key, values := range r.ctx.Request.PostForm {
		if len(values) < 1 {
			continue
		}

		map1[key] = values[0]
	}

	return map1
}

func (r *Request) GetClientIp() string {
	ip := r.GetHeader("X-Forwarded-For")

	if ip == "" {
		ip = r.GetHeader("X-Real-IP")
	}

	if ip == "" {
		ip = r.ctx.ClientIP()
	}

	regex1 := regexp.MustCompile(RegexConst.CommaSep)
	parts := regex1.Split(strings.TrimSpace(ip), -1)

	if len(parts) < 1 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func (r *Request) PathvariableString(name string, defaultValue ...interface{}) string {
	var dv string

	if len(defaultValue) > 0 {
		if s1, err := castx.ToStringE(defaultValue[0]); err == nil {
			dv = s1
		}
	}

	value := r.ctx.Param(name)

	if value == "" {
		return dv
	}

	return value
}

func (r *Request) PathvariableBool(name string, defaultValue ...interface{}) bool {
	var dv bool

	if len(defaultValue) > 0 {
		if b1, err := castx.ToBoolE(defaultValue[0]); err == nil {
			dv = b1
		}
	}

	if b1, err := castx.ToBoolE(r.ctx.Param(name)); err == nil {
		return b1
	}

	return dv
}

func (r *Request) PathvariableInt(name string, defaultValue ...interface{}) int {
	dv := math.MinInt32

	if len(defaultValue) > 0 {
		if n1, err := castx.ToIntE(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	return castx.ToInt(r.ctx.Param(name), dv)
}

func (r *Request) PathvariableInt64(name string, defaultValue ...interface{}) int64 {
	dv := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToInt64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	return castx.ToInt64(r.ctx.Param(name), dv)
}

func (r *Request) PathvariableFloat32(name string, defaultValue ...interface{}) float32 {
	dv := float32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat32E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	return castx.ToFloat32(r.ctx.Param(name), dv)
}

func (r *Request) PathvariableFloat64(name string, defaultValue ...interface{}) float64 {
	dv := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	return castx.ToFloat64(r.ctx.Param(name), dv)
}

func (r *Request) ParamString(name string, defaultValue ...interface{}) string {
	var dv string

	if len(defaultValue) > 0 {
		if s1, err := castx.ToStringE(defaultValue[0]); err == nil {
			dv = s1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	if s1, err := castx.ToStringE(map1[name]); err == nil {
		return s1
	}

	return dv
}

func (r *Request) ParamStringWithSecurityMode(name string, mode int, defaultValue ...interface{}) string {
	var dv string

	if len(defaultValue) > 0 {
		if s1, err := castx.ToStringE(defaultValue[0]); err == nil {
			dv = s1
		}
	}

	modes := []int{0, 1, 2}

	if !slicex.InIntSlice(mode, modes) {
		mode = 2
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	if s1, err := castx.ToStringE(map1[name]); err == nil {
		switch mode {
		case 1, 2:
			s1 = stringx.StripTags(s1)
		}

		return s1
	}

	return dv
}

func (r *Request) ParamBool(name string, defaultValue ...interface{}) bool {
	var dv bool

	if len(defaultValue) > 0 {
		if b1, err := castx.ToBoolE(defaultValue[0]); err == nil {
			dv = b1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	if b1, err := castx.ToBoolE(map1[name]); err == nil {
		return b1
	}

	return dv
}

func (r *Request) ParamInt(name string, defaultValue ...interface{}) int {
	dv := math.MinInt32

	if len(defaultValue) > 0 {
		if n1, err := castx.ToIntE(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	return castx.ToInt(map1[name], dv)
}

func (r *Request) ParamInt64(name string, defaultValue ...interface{}) int64 {
	dv := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToInt64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	return castx.ToInt64(map1[name], dv)
}

func (r *Request) ParamFloat32(name string, defaultValue ...interface{}) float32 {
	dv := float32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat32E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	return castx.ToFloat32(map1[name], dv)
}

func (r *Request) ParamFloat64(name string, defaultValue ...interface{}) float64 {
	dv := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	map1 := map[string]string{}

	for name, value := range r.GetQueryParams() {
		map1[name] = value
	}

	for name, value := range r.GetFormData() {
		map1[name] = value
	}

	return castx.ToFloat64(map1[name], dv)
}

func (r *Request) GetJwt() *jwt.Token {
	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return nil
	}

	tk, _ := ParseJsonWebToken(token)
	return tk
}

func (r *Request) JwtClaimString(name string, defaultValue ...interface{}) string {
	var dv string

	if len(defaultValue) > 0 {
		if s1, err := castx.ToStringE(defaultValue[0]); err == nil {
			dv = s1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimString(token, name, dv)
}

func (r *Request) JwtClaimBool(name string, defaultValue ...interface{}) bool {
	var dv bool

	if len(defaultValue) > 0 {
		if b1, err := castx.ToBoolE(defaultValue[0]); err == nil {
			dv = b1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimBool(token, name, dv)
}

func (r *Request) JwtClaimInt(name string, defaultValue ...interface{}) int {
	dv := math.MinInt32

	if len(defaultValue) > 0 {
		if n1, err := castx.ToIntE(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimInt(token, name, dv)
}

func (r *Request) JwtClaimInt64(name string, defaultValue ...interface{}) int64 {
	dv := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToInt64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimInt64(token, name, dv)
}

func (r *Request) JwtClaimFloat32(name string, defaultValue ...interface{}) float32 {
	dv := float32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat32E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimFloat32(token, name, dv)
}

func (r *Request) JwtClaimFloat64(name string, defaultValue ...interface{}) float64 {
	dv := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		if n1, err := castx.ToFloat64E(defaultValue[0]); err == nil {
			dv = n1
		}
	}

	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return dv
	}

	return JwtClaimFloat64(token, name, dv)
}

func (r *Request) JwtClaimStringSlice(name string) []string {
	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return make([]string, 0)
	}

	return JwtClaimStringSlice(token, name)
}

func (r *Request) JwtClaimIntSlice(name string) []int {
	token := strings.TrimSpace(r.GetHeader("Authorization"))
	token = stringx.RegexReplace(token, `[\x20\t]+`, " ")

	if strings.Contains(token, " ") {
		token = stringx.SubstringAfter(token, " ")
	}

	if token == "" {
		return make([]int, 0)
	}

	return JwtClaimIntSlice(token, name)
}

func (r *Request) GetRawBody() []byte {
	method := r.GetMethod()
	contentType := strings.ToLower(r.GetHeader("Content-Type"))
	isPostForm := strings.Contains(contentType, "application/x-www-form-urlencoded")
	isMultipartForm := strings.Contains(contentType, "multipart/form-data")

	if method == "POST" && (isPostForm || isMultipartForm) {
		formData := r.GetFormData()

		if len(formData) < 1 {
			return make([]byte, 0)
		}

		values := url.Values{}

		for name, value := range formData {
			values[name] = []string{value}
		}

		contents := values.Encode()
		return []byte(contents)
	}

	methods := []string{"POST", "PUT", "PATCH", "DELETE"}

	if !slicex.InStringSlice(method, methods) {
		return make([]byte, 0)
	}

	isJson := strings.Contains(contentType, "application/json")
	isXml1 := strings.Contains(contentType, "application/xml")
	isXml2 := strings.Contains(contentType, "text/xml")

	if !isJson && !isXml1 && !isXml2 {
		return make([]byte, 0)
	}

	rawBody, ok := r.ctx.Get("requestRawBody")

	if !ok {
		return make([]byte, 0)
	}

	if buf, ok := rawBody.([]byte); ok && len(buf) > 0 {
		return buf
	}

	return make([]byte, 0)
}

// @param string[]|string rules
func (r *Request) GetMap(rules ...interface{}) map[string]interface{} {
	method := r.GetMethod()
	methods := []string{"POST", "PUT", "PATCH", "DELETE"}
	contentType := strings.ToLower(r.GetHeader("Content-Type"))
	isPostForm := strings.Contains(contentType, "application/x-www-form-urlencoded")
	isMultipartForm := strings.Contains(contentType, "multipart/form-data")
	isJson := strings.Contains(contentType, "application/json")
	isXml1 := strings.Contains(contentType, "application/xml")
	isXml2 := strings.Contains(contentType, "text/xml")
	map1 := map[string]interface{}{}

	if method == "GET" {
		for key, value := range r.GetQueryParams() {
			map1[key] = value
		}
	} else if method == "POST" && (isPostForm || isMultipartForm) {
		for key, value := range r.GetQueryParams() {
			map1[key] = value
		}

		for key, value := range r.GetFormData() {
			map1[key] = value
		}
	} else if slicex.InStringSlice(method, methods) {
		return map1
	} else if isJson {
		map1 = jsonx.MapFrom(r.GetRawBody())
	} else if isXml1 || isXml2 {
		map2 := mapx.FromXml(r.GetRawBody())

		for key, value := range map2 {
			map1[key] = value
		}
	}

	if len(map1) < 1 {
		return map[string]interface{}{}
	}

	if len(rules) < 1 {
		return map1
	}

	return mapx.FromRequestParam(map1, rules...)
}

func (r *Request) DtoBind(dto interface{}) error {
	return mapx.BindToDto(r.GetMap(), dto)
}

func (r *Request) GetUploadedFile(formFieldName string) *multipart.FileHeader {
	if fh, err := r.ctx.FormFile(formFieldName); err != nil {
		return fh
	}

	return nil
}
