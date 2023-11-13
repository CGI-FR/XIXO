package xixo

import (
	"fmt"
	"strings"
)

type CommentElement struct {
	OuterTextBefore string
	Comment         string
}

func (c CommentElement) String() string {
	return fmt.Sprintf("%s<!--%s-->", c.OuterTextBefore, c.Comment)
}

type XMLElement struct {
	Name      string
	Attrs     map[string]string
	AttrKeys  []string
	InnerText string
	Childs    map[string][]XMLElement
	Err       error

	// filled when xpath enabled
	childs    []*XMLElement
	parent    *XMLElement
	attrs     []*xmlAttr
	localName string
	prefix    string

	outerTextBefore string
	comments        []CommentElement
}

type xmlAttr struct {
	name  string
	value string
}

// SelectElements finds child elements with the specified xpath expression.
func (n *XMLElement) SelectElements(exp string) ([]*XMLElement, error) {
	return find(n, exp)
}

// SelectElement finds child elements with the specified xpath expression.
func (n *XMLElement) SelectElement(exp string) (*XMLElement, error) {
	return findOne(n, exp)
}

func (n *XMLElement) FirstChild() *XMLElement {
	if n.childs == nil {
		return nil
	}

	if len(n.childs) > 0 {
		return n.childs[0]
	}

	return nil
}

func (n *XMLElement) LastChild() *XMLElement {
	if l := len(n.childs); l > 0 {
		return n.childs[l-1]
	}

	return nil
}

func (n *XMLElement) PrevSibling() *XMLElement {
	if n.parent != nil {
		for i, c := range n.parent.childs {
			if c == n {
				if i >= 0 {
					return n.parent.childs[i-1]
				}

				return nil
			}
		}
	}

	return nil
}

func (n *XMLElement) NextSibling() *XMLElement {
	if n.parent != nil {
		for i, c := range n.parent.childs {
			if c == n {
				if i+1 < len(n.parent.childs) {
					return n.parent.childs[i+1]
				}

				return nil
			}
		}
	}

	return nil
}

func (n *XMLElement) String() string {
	xmlChilds := ""

	for node := n.FirstChild(); node != nil; node = node.NextSibling() {
		xmlChilds += node.String()
	}

	attributes := n.Name + " "
	for _, key := range n.AttrKeys {
		attributes += fmt.Sprintf("%s=\"%s\" ", key, n.Attrs[key])
	}

	attributes = strings.Trim(attributes, " ")

	commentsString := ""
	for _, comment := range n.comments {
		commentsString += comment.String()
	}

	return fmt.Sprintf("%s<%s>%s%s%s</%s>",
		n.outerTextBefore,
		attributes,
		commentsString,
		xmlChilds,
		n.InnerText,
		n.Name)
}

func (n *XMLElement) AddAttribute(name string, value string) {
	if n.Attrs == nil {
		n.Attrs = make(map[string]string)
	}
	// if name don't exsite in Attrs yet
	if _, ok := n.Attrs[name]; !ok {
		// Add un key in slice to keep the order of attributes
		n.AttrKeys = append(n.AttrKeys, name)
	}
	// change the value of attribute
	n.Attrs[name] = value
}

func (n *XMLElement) AddComment(comment CommentElement) {
	n.comments = append(n.comments, comment)
}

func NewXMLElement() *XMLElement {
	return &XMLElement{
		Name:      "",
		Attrs:     map[string]string{},
		AttrKeys:  make([]string, 0),
		InnerText: "",
		Childs:    map[string][]XMLElement{},
		Err:       nil,
		childs:    []*XMLElement{},
		parent:    nil,
		attrs:     []*xmlAttr{},
		localName: "",
		prefix:    "",
	}
}
