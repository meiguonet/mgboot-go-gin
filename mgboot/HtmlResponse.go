package mgboot

type HtmlResponse struct {
	contents string
}

func NewHtmlResponse(contents string) HtmlResponse {
	return HtmlResponse{contents: contents}
}

func (p HtmlResponse) GetContentType() string {
	return "text/html; charset=utf-8"
}

func (p HtmlResponse) GetContents() (int, string) {
	return 200, p.contents
}
