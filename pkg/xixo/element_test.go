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

	attr := xixo.Attribute{"foo", "bar", xixo.SimpleQuote}

	root.AddAttribute(attr)

	expected := map[string]xixo.Attribute{"foo": attr}

	assert.Equal(t, root.Attrs, expected)
}

func TestAddAttributsShouldInOutputWithString(t *testing.T) {
	t.Parallel()

	root := xixo.NewXMLElement()
	root.Name = parentTag
	root.InnerText = "Hello"
	root.AddAttribute(xixo.Attribute{"foo", "bar", xixo.DoubleQuotes})

	expected := "<root foo=\"bar\">Hello</root>"
	assert.Equal(t, expected, root.String())
}

func TestEditAttributsShouldInOutputWithString(t *testing.T) {
	t.Parallel()

	root := xixo.NewXMLElement()
	root.Name = parentTag
	root.InnerText = "Hello"
	root.AddAttribute(xixo.Attribute{"foo", "bar", xixo.DoubleQuotes})

	expected := "<root foo=\"bar\">Hello</root>"
	assert.Equal(t, expected, root.String())
	root.AddAttribute(xixo.Attribute{"foo", "bas", xixo.DoubleQuotes})

	expected = "<root foo=\"bas\">Hello</root>"
	assert.Equal(t, expected, root.String())
}

func TestElementStringShouldRemoveTargetAttribute(t *testing.T) {
	t.Parallel()

	rootXML := `<root location="Nantes">
  <element1 name="joe" age="5">Hello world !</element1>
  <element2 name="doe">Contenu2 </element2>
</root>`

	var resultXMLBuffer bytes.Buffer
	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), &resultXMLBuffer).EnableXpath()
	parser.RegisterMapCallback(parentTag, func(x map[string]string) (map[string]string, error) {
		delete(x, "element1@name")
		delete(x, "@location")
		delete(x, "element2@name")

		return x, nil
	})

	expect := `<root>
  <element1 age="5">Hello world !</element1>
  <element2>Contenu2 </element2>
</root>`

	err := parser.Stream()
	assert.Nil(t, err)

	resultXML := resultXMLBuffer.String()
	assert.Equal(t, expect, resultXML)
}

func TestElementStringShouldRemoveTargetTag(t *testing.T) {
	t.Parallel()

	rootXML := `<root location="Nantes" name="Agency">
  <element1 name="joe" age="5">Hello world !</element1>
<element2 name="doe">Contenu2 </element2>
</root>`

	var resultXMLBuffer bytes.Buffer
	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), &resultXMLBuffer).EnableXpath()
	parser.RegisterMapCallback(parentTag, func(x map[string]string) (map[string]string, error) {
		delete(x, "element1")
		delete(x, "@location")

		return x, nil
	})

	expect := `<root name="Agency">
<element2 name="doe">Contenu2 </element2>
</root>`

	err := parser.Stream()
	assert.Nil(t, err)

	resultXML := resultXMLBuffer.String()
	assert.Equal(t, expect, resultXML)
}
