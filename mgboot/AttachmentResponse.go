package mgboot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meiguonet/mgboot-go-common/util/mimex"
	"io/ioutil"
)

type AttachmentResponse struct {
	buf                []byte
	mimeType           string
	attachmentFileName string
}

func NewAttachmentResponseFromFile(fpath, attachmentFileName string, mimeType ...string) AttachmentResponse {
	buf, _ := ioutil.ReadFile(fpath)
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return AttachmentResponse{
		buf:                buf,
		mimeType:           _mimeType,
		attachmentFileName: attachmentFileName,
	}
}

func NewAttachmentResponseFromBuffer(buf []byte, attachmentFileName string, mimeType ...string) AttachmentResponse {
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return AttachmentResponse{
		buf:                buf,
		mimeType:           _mimeType,
		attachmentFileName: attachmentFileName,
	}
}

func (p AttachmentResponse) GetContentType() string {
	if p.mimeType == "" {
		return "application/octet-stream"
	}

	return p.mimeType
}

func (p AttachmentResponse) GetContents() (int, string) {
	if len(p.buf) < 1 || p.attachmentFileName == "" {
		return 400, ""
	}
	
	return 200, ""
}

func (p AttachmentResponse) Buffer() []byte {
	return p.buf
}

func (p AttachmentResponse) AddSpecifyHeaders(ctx *gin.Context) {
	disposition := fmt.Sprintf(`attachment; filename="%s"`, p.attachmentFileName)
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(p.buf)))
	ctx.Header("Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", disposition)
	ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	ctx.Header("Pragma", "public")
}
