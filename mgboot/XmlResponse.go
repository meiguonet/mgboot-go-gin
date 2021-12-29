package mgboot

type XmlResponse struct {
	contents string
}

func NewXmlResponse(contents string) XmlResponse {
	return XmlResponse{contents: contents}
}

func (p XmlResponse) GetContentType() string {
	return "text/xml; charset=utf-8"
}

func (p XmlResponse) GetContents() (int, string) {
	return 200, p.contents
}
