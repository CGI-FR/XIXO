package xixo_test

import (
	"encoding/json"
	"testing"

	"github.com/CGI-FR/xixo/pkg/xixo"
	"github.com/stretchr/testify/assert"
)

const (
	newChildContent = "newChildContent"
	rootXML         = `<root>
	<element1>Hello world !</element1>
	<element2>Contenu2 </element2>
</root>`
)

func mapCallback(dict map[string]string) (map[string]string, error) {
	dict["element1"] = newChildContent

	return dict, nil
}

// TestMapCallback should convert.
func TestMapCallback(t *testing.T) {
	t.Parallel()

	element1 := createTreeFromXMLString(rootXML)

	assert.Equal(t, rootXML, element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallback)(element1)
	assert.Nil(t, err)

	text := editedElement1.FirstChild().InnerText

	assert.Equal(t, newChildContent, text)

	expected := `<root>
	<element1>newChildContent</element1>
	<element2>Contenu2 </element2>
</root>`

	assert.Equal(t, expected, editedElement1.String())
}

func jsonCallback(source string) (string, error) {
	dict := map[string]string{}

	err := json.Unmarshal([]byte(source), &dict)
	if err != nil {
		return "", err
	}

	dict["element1"] = newChildContent

	result, err := json.Marshal(dict)

	return string(result), err
}

// TestMapCallback should convert.
func TestJsonCallback(t *testing.T) {
	t.Parallel()

	root := createTreeFromXMLString(rootXML)

	editedRoot, err := xixo.XMLElementToJSONCallback(jsonCallback)(root)
	assert.Nil(t, err)

	assert.Equal(t,
		"<root>\n\t<element1>newChildContent</element1>\n\t<element2>Contenu2 </element2>\n</root>",
		editedRoot.String(),
	)
}

func badJSONCallback(source string) (string, error) {
	return "{ hello: 1 " + source, nil
}

func TestBadJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createTreeFromXMLString(rootXML)

	_, err := xixo.XMLElementToJSONCallback(badJSONCallback)(element1)

	assert.NotNil(t, err)
}

func TestMapCallbackWithAttributs(t *testing.T) {
	t.Parallel()

	rootXML := `<root>
		<element1 age="22">Hello world !</element1>
		<element2>Contenu2 </element2>
	</root>`

	element1 := createTreeFromXMLString(rootXML)

	assert.Equal(t, rootXML, element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallbackAttributs)(element1)
	assert.Nil(t, err)

	text := editedElement1.FirstChild().InnerText

	assert.Equal(t, "newChildContent", text)

	expected := `<root>
		<element1 age="50">newChildContent</element1>
		<element2>Contenu2 </element2>
	</root>`

	assert.Equal(t, expected, editedElement1.String())
}

func mapCallbackAttributs(dict map[string]string) (map[string]string, error) {
	dict["element1@age"] = "50"
	dict["element1"] = newChildContent

	return dict, nil
}

func TestMapCallbackWithAttributsParentAndChilds(t *testing.T) {
	t.Parallel()

	rootXML := `<root type="foo">
		<element1 age="22" sex="male">Hello world !</element1>
		<element2>Contenu2 </element2>
	</root>`

	element1 := createTreeFromXMLString(rootXML)

	assert.Equal(t, rootXML, element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallbackAttributsWithParent)(element1)
	assert.Nil(t, err)

	text := editedElement1.FirstChild().InnerText

	assert.Equal(t, "newChildContent", text)

	expected := `<root type="bar">
		<element1 age="50" sex="male">newChildContent</element1>
		<element2 age="25">Contenu2 </element2>
	</root>`

	assert.Equal(t, expected, editedElement1.String())
}

func mapCallbackAttributsWithParent(dict map[string]string) (map[string]string, error) {
	dict["@type"] = "bar"
	dict["element1@age"] = "50"
	dict["element1"] = newChildContent
	dict["element2@age"] = "25"

	return dict, nil
}
