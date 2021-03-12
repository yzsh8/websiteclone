package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

/**
* 获取html文档的head
 */
func getHead(res string) (*html.Node, error) {
	doc, _ := html.Parse(strings.NewReader(res))

	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "head" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <head> in the node tree")
}

/**
* 获取html文档的body
 */
func getBody(res string) (*html.Node, error) {
	doc, _ := html.Parse(strings.NewReader(res))

	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

/**
* 读取指定内容
 */
func getNodeById(res string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(res))

	if err != nil {
		fmt.Println("加载节点失败", err.Error())
	} else {
		fmt.Println("type->", doc.Type, "data->", doc.Data, "attr->", doc.Attr, "namespace->", doc.Namespace)
	}

	var b *html.Node
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		b = c
		//fmt.Println("debug->", c.Attr)
		if c.Data == "div" {
			fmt.Println("debug run", "div")
		}
	}

	return b, nil
}

/**
* 读取内容
 */
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
