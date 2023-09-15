package xixo

import (
	"fmt"

	"github.com/tamerh/xpath"
)

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func (x *XMLParser) CreateXPathNavigator(top *XMLElement) *XMLNodeNavigator {
	return &XMLNodeNavigator{curr: top, root: top, attr: -1}
}

// Compile the given xpath expression.
func (x *XMLParser) CompileXpath(expr string) (*xpath.Expr, error) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func createXPathNavigator(top *XMLElement) *XMLNodeNavigator {
	return &XMLNodeNavigator{curr: top, root: top, attr: -1}
}

type XMLNodeNavigator struct {
	root, curr *XMLElement
	attr       int
}

// Find searches the Node that matches by the specified XPath expr.
func find(top *XMLElement, expr string) ([]*XMLElement, error) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		return []*XMLElement{}, err
	}

	t := exp.Select(createXPathNavigator(top))

	var elems []*XMLElement

	for t.MoveNext() {
		current, ok := t.Current().(*XMLNodeNavigator)
		if !ok {
			return nil, fmt.Errorf("current is not a XMLNodeNavigator %v", current)
		}

		elems = append(elems, current.curr)
	}

	return elems, nil
}

// FindOne searches the Node that matches by the specified XPath expr,
// and returns first element of matched.
func findOne(top *XMLElement, expr string) (*XMLElement, error) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}

	t := exp.Select(createXPathNavigator(top))

	var elem *XMLElement

	if t.MoveNext() {
		navigator, ok := t.Current().(*XMLNodeNavigator)
		if !ok {
			return nil, fmt.Errorf("Current is not a XMLNodeNavigator %v", navigator)
		}

		elem = navigator.curr
	}

	return elem, nil
}

func (x *XMLNodeNavigator) Current() *XMLElement {
	return x.curr
}

func (x *XMLNodeNavigator) NodeType() xpath.NodeType {
	if x.curr == x.root {
		return xpath.RootNode
	}

	if x.attr != -1 {
		return xpath.AttributeNode
	}

	return xpath.ElementNode
}

func (x *XMLNodeNavigator) LocalName() string {
	if x.attr != -1 {
		return x.curr.attrs[x.attr].name
	}

	return x.curr.localName
}

func (x *XMLNodeNavigator) Prefix() string {
	return x.curr.prefix
}

func (x *XMLNodeNavigator) Value() string {
	if x.attr != -1 {
		return x.curr.attrs[x.attr].value
	}

	return x.curr.InnerText
}

func (x *XMLNodeNavigator) Copy() xpath.NodeNavigator {
	n := *x

	return &n
}

func (x *XMLNodeNavigator) MoveToRoot() {
	x.curr = x.root
}

func (x *XMLNodeNavigator) MoveToParent() bool {
	if x.attr != -1 {
		x.attr = -1

		return true
	} else if node := x.curr.parent; node != nil {
		x.curr = node

		return true
	}

	return false
}

func (x *XMLNodeNavigator) MoveToNextAttribute() bool {
	if x.attr >= len(x.curr.attrs)-1 {
		return false
	}
	x.attr++

	return true
}

func (x *XMLNodeNavigator) MoveToChild() bool {
	if node := x.curr.FirstChild(); node != nil {
		x.curr = node

		return true
	}

	return false
}

func (x *XMLNodeNavigator) MoveToFirst() bool {
	if x.curr.parent != nil {
		node := x.curr.parent.FirstChild()
		if node != nil {
			x.curr = node

			return true
		}
	}

	return false
}

func (x *XMLNodeNavigator) MoveToPrevious() bool {
	node := x.curr.PrevSibling()
	if node != nil {
		x.curr = node

		return true
	}

	return false
}

func (x *XMLNodeNavigator) MoveToNext() bool {
	node := x.curr.NextSibling()
	if node != nil {
		x.curr = node

		return true
	}

	return false
}

func (x *XMLNodeNavigator) String() string {
	return x.Value()
}

func (x *XMLNodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*XMLNodeNavigator)
	if !ok || node.root != x.root {
		return false
	}

	x.curr = node.curr
	x.attr = node.attr

	return true
}
