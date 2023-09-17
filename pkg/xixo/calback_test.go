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
	assert.Equal(t, "<element1><child1></child1></element1>", element1.String())

	editedElement1, err := xixo.XMLElementToMapCallback(mapCallback)(&element1)
	assert.Nil(t, err)

	text := editedElement1.Childs["child1"][0].InnerText

	assert.Equal(t, "newChildContent", text)

	assert.Equal(t, "<element1><child1>newChildContent</child1></element1>", editedElement1.String())
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

func badJSONCallback(source string) (string, error) {
	return "{ hello: 1 " + source, nil
}

func TestBadJsonCallback(t *testing.T) {
	t.Parallel()

	element1 := createElement1()

	_, err := xixo.XMLElementToJSONCallback(badJSONCallback)(&element1)

	assert.NotNil(t, err)
}
