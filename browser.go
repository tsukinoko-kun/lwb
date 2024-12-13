package lwb

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"

	"github.com/tsukinoko-kun/lwb/util"
	"golang.org/x/net/html"
)

type Browser struct {
	url       string
	userAgent string
	http      *http.Client
	document  *html.Node
	cookies   *cookiejar.Jar
	mut       sync.RWMutex
}

func NewBrowser(userAgent string) (*Browser, error) {
	cj, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{
		Jar: cj,
	}

	b := &Browser{
		userAgent: userAgent,
		cookies:   cj,
		http:      hc,
	}

	return b, nil
}

func (b *Browser) Get(url string) error {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.url = url

	resp, err := b.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the HTML document
	b.document, err = html.Parse(resp.Body)
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
