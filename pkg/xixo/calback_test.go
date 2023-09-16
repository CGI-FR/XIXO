package xixo_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func mapCallback(dict map[string]string) map[string]string {
	dict["child1"] = "newChildContent"
	return dict
}

// TestMapCallback should convert
func TestMapCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	editedElement1 := xixo.XMLElementToMapCallback(mapCallback)(&element1)

	text := editedElement1.Childs["child1"][0].InnerText

	assert.Equal(t, "newChildContent", text)
}

func jsonCallback(source string) string {
	dict := map[string]string{}

	json.Unmarshal([]byte(source), &dict)
	dict["child1"] = "newChildContent"

	result, _ := json.Marshal(dict)

	return string(result)
}

// TestMapCallback should convert
func TestJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	editedElement1 := xixo.XMLElementToJSONCallback(jsonCallback)(&element1)

	text := editedElement1.Childs["child1"][0].InnerText

	assert.Equal(t, "newChildContent", text)
}

func createElement1() xixo.XMLElement {
	//nolint
	element1 := xixo.XMLElement{}

	element1.Name = "element1"
	element1.Childs = map[string][]xixo.XMLElement{}

	//nolint
	child1 := xixo.XMLElement{}

	child1.Name = "child1"

	element1.Childs[child1.Name] = []xixo.XMLElement{child1}

	return element1
}
