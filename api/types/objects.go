package types

type Object struct {
	Name          string
	ContentType   string
	ContentLength string
}

type Objects struct {
	Objects []Object `xml:"Object"`
}
