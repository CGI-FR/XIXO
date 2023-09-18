package xixo_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func createElement1() *xixo.XMLElement {
	rootXML := `
	<root>
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
	if err != nil {
		return nil
	}

	return root
}

func TestElementStringShouldReturnXML(t *testing.T) {
	t.Parallel()

	rootXML := `
	<root>
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
  <element3>Contenu2 </element3>
  <element4>Contenu2 </element4>
  <element5>Contenu2 </element5>
</root>`

	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard).EnableXpath()
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x

		return x, nil
	})

	err := parser.Stream()
	assert.Nil(t, err)

	assert.Equal(t, rootXML, root.String())
}
