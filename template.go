package gozt

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

type Template struct {
	Document *goquery.Document
}

func ReadTemplate(r io.Reader) (*Template, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return &Template{Document: doc}, nil
}
