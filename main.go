package lwb

import (
	"net/http"
	"strings"
	"sync"

	"github.com/tsukinoko-kun/lwb/util"
	"golang.org/x/net/html"
)

type (
	Browser struct {
		url       string
		userAgent string
		document  *html.Node
		mut       sync.RWMutex
	}

	Element struct {
		node *html.Node
	}
)

func NewBrowser(userAgent string) *Browser {
	b := &Browser{
		userAgent: userAgent,
	}

	return b
}

func (b *Browser) Get(url string) error {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.url = url

	rest, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rest.Body.Close()

	b.document, err = html.Parse(rest.Body)
	if err != nil {
		return err
	}

	return nil
}

func (b *Browser) GetElementById(id string) *Element {
	b.mut.RLock()
	defer b.mut.RUnlock()

	var nodes util.Stack[*html.Node] = []*html.Node{b.document}
	for !nodes.Empty() {
		node := nodes.Pop()
		for _, a := range node.Attr {
			if strings.ToLower(a.Key) != "id" {
				continue
			}
			if a.Val != id {
				break
			}
			return &Element{node: node}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodes.Push(c)
		}
	}
	return nil
}

func (b *Browser) GetElementsByClassName(class string) []*Element {
	b.mut.RLock()
	defer b.mut.RUnlock()

	var elements []*Element
	var nodes util.Stack[*html.Node] = []*html.Node{b.document}
	for !nodes.Empty() {
		node := nodes.Pop()
	attr_loop:
		for _, a := range node.Attr {
			if strings.ToLower(a.Key) != "class" {
				continue
			}
			for _, c := range classNames(a.Val) {
				if c == class {
					elements = append(elements, &Element{node: node})
					break attr_loop
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodes.Push(c)
		}
	}
	return nil
}

func classNames(class string) []string {
	return strings.Split(class, " ")
}
