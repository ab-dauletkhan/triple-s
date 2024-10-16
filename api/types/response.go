package types

type Error struct {
	Code    int    `xml:"Code"`
	Message string `xml:"Message"`
}
