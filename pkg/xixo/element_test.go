package xixo_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/CGI-FR/xixo/pkg/xixo"
	"github.com/stretchr/testify/assert"
)

const (
	parentTag = "root"
)

func createTreeFromXMLString(rootXML string) *xixo.XMLElement {
	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x

		return x, nil
	})

	err := parser.Stream()
	if err != nil {
		return nil
	}

	return root
}

func TestElementStringShouldReturnXML(t *testing.T) {
	t.Parallel()

	rootXML := `<root>
		<element1>Hello world !</element1>
		<element2>Contenu2 </element2>
	</root>`

	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x

		return x, nil
	})

	err := parser.Stream()
	assert.Nil(t, err)

	expected := `<root>
		<element1>Hello world !</element1>
		<element2>Contenu2 </element2>
	</root>`

	assert.Equal(t, expected, root.String())
}

func TestElementStringShouldReturnXMLWithSameOrder(t *testing.T) {
	t.Parallel()

	rootXML := `<root>
  <element1>Hello world !</element1>
  <element2>Contenu2 </element2>
  <element3>Contenu3 </element3>
  <element4>Contenu4 </element4>
  <element5>Contenu5 </element5>
</root>`

	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x
		assert.Equal(t, root.InnerText, "\n")

		return x, nil
	})

	err := parser.Stream()
	assert.Nil(t, err)

	assert.Equal(t, rootXML, root.String())
}

func TestElementStringShouldPreserverContentOrder(t *testing.T) {
	t.Parallel()

	rootXML := `<root>
  <element1>Hello world !</element1>
  <element2>Contenu2 </element2>
  <element2>Contenu3 </element2>
  <element2>Contenu4 </element2>
  <element2>Contenu5 </element2>
</root>`

	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x
		assert.Equal(t, root.InnerText, "\n")

		return x, nil
	})

	err := parser.Stream()
	assert.Nil(t, err)

	assert.Equal(t, rootXML, root.String())
}

func TestCreatNewXMLElement(t *testing.T) {
	t.Parallel()

	rootXML := `
	<root>
	</root>`

	var root *xixo.XMLElement

	name := parentTag
	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	root = xixo.NewXMLElement()
	root.Name = name
	err := parser.Stream()

	assert.Nil(t, err)

	expected := `<root></root>`

	assert.Equal(t, expected, root.String())
}

func TestAddAttributsShouldSaved(t *testing.T) {
	t.Parallel()

	var root *xixo.XMLElement

	name := parentTag
	root = xixo.NewXMLElement()
	root.Name = name

	root.AddAttribute("foo", "bar")

	expected := map[string]string{"foo": "bar"}

	assert.Equal(t, root.Attrs, expected)
}

func TestAddAttributsShouldInOutputWithString(t *testing.T) {
	t.Parallel()

	root := xixo.NewXMLElement()
	root.Name = parentTag
	root.InnerText = "Hello"
	root.AddAttribute("foo", "bar")

	expected := "<root foo=\"bar\">Hello</root>"
	assert.Equal(t, expected, root.String())
}

func TestEditAttributsShouldInOutputWithString(t *testing.T) {
	t.Parallel()

	root := xixo.NewXMLElement()
	root.Name = parentTag
	root.InnerText = "Hello"
	root.AddAttribute("foo", "bar")

	expected := "<root foo=\"bar\">Hello</root>"
	assert.Equal(t, expected, root.String())
	root.AddAttribute("foo", "bas")

	expected = "<root foo=\"bas\">Hello</root>"
	assert.Equal(t, expected, root.String())
}
