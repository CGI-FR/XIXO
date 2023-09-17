package xixo_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func createElement1() xixo.XMLElement {
	element1 := xixo.XMLElement{}

	element1.Name = "element1"
	element1.Childs = map[string][]xixo.XMLElement{}

	child1 := xixo.XMLElement{}

	child1.Name = "child1"

	element1.Childs[child1.Name] = []xixo.XMLElement{child1}

	return element1
}

func TestElementStringShouldReturnXML(t *testing.T) {
	t.Parallel()

	rootXML := `
	<root>
		<element1>Hello world !</element1>
		<element2>Contenu2 </element2>
	</root>`

	var root *xixo.XMLElement

	parser := xixo.NewXMLParser(bytes.NewBufferString(rootXML), io.Discard)
	parser.RegisterCallback("root", func(x *xixo.XMLElement) (*xixo.XMLElement, error) {
		root = x

		return x, nil
	})

	err := parser.Stream()
	assert.Nil(t, err)

	assert.Equal(t, "<root><element1>Hello world !</element1><element2>Contenu2 </element2></root>", root.String())
}
