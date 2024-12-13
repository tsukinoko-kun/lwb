package lwb

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

type Element struct {
	node    *html.Node
	browser *Browser
}

var ErrorNotClickable = errors.New("element is not clickable")

func (self *Element) Click() error {
	for e := self.node; e != nil; e = e.Parent {
		if e.Type != html.ElementNode {
			continue
		}
		if strings.ToLower(e.Data) == "a" {
			for _, a := range e.Attr {
				if a.Key == "href" {
					return self.browser.Get(a.Val)
				}
			}
		}
	}

	return ErrorNotClickable
}
