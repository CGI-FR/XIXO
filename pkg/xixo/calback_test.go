package xixo_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func mapCallback(dict map[string]string) (map[string]string, error) {
	dict["element1"] = "newChildContent"

	return dict, nil
}

// TestMapCallback should convert.
func TestMapCallback(t *testing.T) {
	t.Parallel()

	element1 := createTree()
	//nolint
	assert.Equal(t, "<root>\n  <element1>Hello world !</element1>\n  <element2>Contenu2 </element2>\n</root>", element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallback)(element1)
	assert.Nil(t, err)

	text := editedElement1.FirstChild().InnerText

	assert.Equal(t, "newChildContent", text)

	//nolint
	assert.Equal(t, "<root>\n  <element1>newChildContent</element1>\n  <element2>Contenu2 </element2>\n</root>", editedElement1.String())
}

func jsonCallback(source string) (string, error) {
	dict := map[string]string{}

	err := json.Unmarshal([]byte(source), &dict)
	if err != nil {
		return "", err
	}

	dict["element1"] = "newChildContent"

	result, err := json.Marshal(dict)

	return string(result), err
}

// TestMapCallback should convert.
func TestJsonCallback(t *testing.T) {
	t.Parallel()

	root := createTree()

	editedRoot, err := xixo.XMLElementToJSONCallback(jsonCallback)(root)
	assert.Nil(t, err)

	element1, err := editedRoot.SelectElement("element1")

	assert.Nil(t, err)

	assert.Equal(t, "newChildContent", element1.InnerText)
}

func badJSONCallback(source string) (string, error) {
	return "{ hello: 1 " + source, nil
}

func TestBadJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createTree()

	_, err := xixo.XMLElementToJSONCallback(badJSONCallback)(element1)

	assert.NotNil(t, err)
}

func TestMapCallbackWithAttributs(t *testing.T) {
	t.Parallel()

	element1 := createTree()
	//nolint
	assert.Equal(t, "<root>\n  <element1 age='22'>Hello world !</element1>\n  <element2>Contenu2 </element2>\n</root>", element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallbackAttributs)(element1)
	assert.Nil(t, err)

	text := editedElement1.FirstChild().InnerText

	assert.Equal(t, "newChildContent", text)

	//nolint
	assert.Equal(t, "<root>\n  <element1 age='50'>newChildContent</element1>\n  <element2>Contenu2 </element2>\n</root>", editedElement1.String())
}

func mapCallbackAttributs(dict map[string]string) (map[string]string, error) {
	dict["element1@age"] = "50"
	dict["element1"] = "newChildContent"

	return dict, nil
}
