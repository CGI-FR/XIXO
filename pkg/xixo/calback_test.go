package xixo_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func mapCallback(dict map[string]string) (map[string]string, error) {
	dict["child1"] = "newChildContent"

	return dict, nil
}

// TestMapCallback should convert.
func TestMapCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallback)(&element1)
	assert.Nil(t, err)

	text := editedElement1.Childs["child1"][0].InnerText

	assert.Equal(t, "newChildContent", text)
}

func jsonCallback(source string) (string, error) {
	dict := map[string]string{}

	err := json.Unmarshal([]byte(source), &dict)
	if err != nil {
		return "", err
	}

	dict["child1"] = "newChildContent"

	result, err := json.Marshal(dict)

	return string(result), err
}

// TestMapCallback should convert.
func TestJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	editedElement1, err := xixo.XMLElementToJSONCallback(jsonCallback)(&element1)
	assert.Nil(t, err)

	text := editedElement1.Childs["child1"][0].InnerText

	assert.Equal(t, "newChildContent", text)
}

func createElement1() xixo.XMLElement {
	element1 := xixo.XMLElement{}

	element1.Name = "element1"
	element1.Childs = map[string][]xixo.XMLElement{}

	child1 := xixo.XMLElement{}

	child1.Name = "child1"

	element1.Childs[child1.Name] = []xixo.XMLElement{child1}

	return element1
}

func badJSONCallback(source string) (string, error) {
	return "{ hello: 1 " + source, nil
}

func TestBadJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	_, err := xixo.XMLElementToJSONCallback(badJSONCallback)(&element1)

	assert.NotNil(t, err)
}
