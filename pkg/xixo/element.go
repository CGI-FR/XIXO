package xixo

import (
	"fmt"
	"strings"
)

type Quote string

const (
	SimpleQuote  = "'"
	DoubleQuotes = "\""
)

func ParseQuoteType(car byte) Quote {
	if car == '\'' {
		return SimpleQuote
	}

	return DoubleQuotes
}

type Attribute struct {
	Name  string
	Value string
	Quote Quote
}

func (attr Attribute) String() string {
	if attr.Quote == SimpleQuote {
		return fmt.Sprintf("%s='%s'", attr.Name, attr.Value)
	}

	return fmt.Sprintf("%s=\"%s\"", attr.Name, attr.Value)
}

type XMLElement struct {
	Name      string
	Attrs     map[string]Attribute
	AttrKeys  []string
	InnerText string
	Childs    map[string][]XMLElement
	Err       error

	// filled when xpath enabled
	childs    []*XMLElement
	parent    *XMLElement
	localName string
	prefix    string

	outerTextBefore string
	autoClosable    bool
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
		attributes += n.Attrs[key].String() + " "
	}

	attributes = strings.Trim(attributes, " ")

	if n.autoClosable && n.InnerText == "" && xmlChilds == "" {
		return fmt.Sprintf("%s<%s/>",
			n.outerTextBefore,
			attributes)
	}

	return fmt.Sprintf("%s<%s>%s%s</%s>",
		n.outerTextBefore,
		attributes,
		xmlChilds,
		n.InnerText,
		n.Name)
}

func (n *XMLElement) AddAttribute(attr Attribute) {
	if n.Attrs == nil {
		n.Attrs = make(map[string]Attribute)
	}
	// if name don't exsite in Attrs yet
	if _, ok := n.Attrs[attr.Name]; !ok {
		// Add un key in slice to keep the order of attributes
		n.AttrKeys = append(n.AttrKeys, attr.Name)
	} else {
		attr.Quote = n.Attrs[attr.Name].Quote
	}
	// change the value of attribute
	n.Attrs[attr.Name] = attr
}

func NewXMLElement() *XMLElement {
	return &XMLElement{
		Name:      "",
		Attrs:     map[string]Attribute{},
		AttrKeys:  make([]string, 0),
		InnerText: "",
		Childs:    map[string][]XMLElement{},
		Err:       nil,
		childs:    []*XMLElement{},
		parent:    nil,
		localName: "",
		prefix:    "",
	}
}
