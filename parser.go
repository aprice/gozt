package gozt

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
	Data     interface{}
	Template *Template

	document  *goquery.Document
	buffer    *bytes.Buffer
	rootScope *scope
}

func New(template *Template, data interface{}) *Parser {
	return &Parser{
		Data:      data,
		Template:  template,
		document:  goquery.CloneDocument(template.Document),
		rootScope: newRootScope(data),
	}
}

func (p *Parser) Read(b []byte) (int, error) {
	if p.buffer == nil || p.buffer.Len() == 0 {
		html, err := p.document.Html()
		if err != nil {
			return 0, err
		}
		p.buffer = bytes.NewBufferString(html)
	}
	return p.buffer.Read(b)
}

func (p *Parser) Parse() error {
	err := p.derive()
	if err != nil {
		return err
	}
	err = p.handleConditionals(p.document.Selection, p.rootScope)
	if err != nil {
		return err
	}
	err = p.handleSnippets()
	if err != nil {
		return err
	}
	err = p.inject()
	if err != nil {
		return err
	}
	html, err := p.document.Html()
	if err != nil {
		return err
	}
	html, err = p.rootScope.substitute(html)
	if err != nil {
		return err
	}
	p.buffer = bytes.NewBufferString(html)
	return nil
}

func (p *Parser) derive() error {
	// Check for data-z-derivesfrom or class,
	// find a matching parent template,
	// partially parse the parent (derivation only),
	// handle overrides,
	// handle appends,
	// then replace our document with the derivation result.
	return nil
}

func (p *Parser) inject() error {
	// Starting at root element,
	// for each child element,
	// look for data-z-inject,
	// or if not found, see if id matches a property in scope,
	// or if not, look for first class that matches a property in scope,
	// and if found,
	// set new scope,
	// 		if it's an array or slice, repeat the element for each item,
	// inject the element,
	// then do substitution in scope.
	return nil
}

func (p *Parser) handleInjection(el *goquery.Selection, scope *scope) error {
	children := el.Children()
	for child := children.First(); child != nil; child = children.Next() {
		// Match attribute or class to model
	}
	return nil
}

func (p *Parser) handleConditionals(el *goquery.Selection, scope *scope) error {
	// For each child element,
	// look for data-z-lorem/data-z-if,
	// and possibly remove the element,
	// otherwise recurse.
	return nil
}

func (p *Parser) handleSnippets() error {
	// Starting at root element,
	// for each child element,
	// check for data-z-snippet,
	// look up the snippet document,
	// partially parse it (derivation only),
	// find the snippet element,
	// and replace our element with it.
	return nil
}
