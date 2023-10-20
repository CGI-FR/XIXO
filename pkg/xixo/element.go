package xixo

import (
	"fmt"
	"strings"
)

type XMLElement struct {
	Name      string
	Attrs     map[string]string
	InnerText string
	Childs    map[string][]XMLElement
	Err       error
	// filled when xpath enabled
	childs    []*XMLElement
	parent    *XMLElement
	attrs     []*xmlAttr
	localName string
	prefix    string
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
		xmlChilds += "  " + node.String() + "\n"
	}

	if len(xmlChilds) > 0 {
		xmlChilds = "\n" + xmlChilds
	}

	attributes := n.Name + " "

	for _, attr := range n.attrs {
		attributes += fmt.Sprintf("%s=\"%s\" ", attr.name, attr.value)
	}

	attributes = strings.Trim(attributes, " ")

	return fmt.Sprintf("<%s>%s%s</%s>",
		attributes,
		n.InnerText,
		xmlChilds,
		n.Name)
}

func (n *XMLElement) AddAttribut(name string, value string) {
	if n.Attrs == nil {
		n.Attrs = make(map[string]string)
	}
	n.Attrs[name] = value
}

func NewXMLElement(name string, attrs map[string]string) *XMLElement {
	return &XMLElement{
		Name:      name,
		Attrs:     attrs,
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
