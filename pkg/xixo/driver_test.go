package xixo_test

import (
	"bytes"
	"testing"

	"github.com/CGI-FR/xixo/pkg/xixo"
	"github.com/stretchr/testify/assert"
)

// TestFuncDriverEdit tests the FuncDriver by creating a reader, writer, and a callback function,
// then it processes the XML and asserts the expected output and that the callback was called.
func TestFuncDriverEdit(t *testing.T) {
	t.Parallel()

	// Create a reader with an XML string, an empty writer, a callback function, and a driver.
	reader := bytes.NewBufferString("<root><element1>innerTexta</element1></root>")
	writer := bytes.Buffer{}
	called := false
	callback := func(input map[string]string) (map[string]string, error) {
		called = true
		input["element1"] = "innerTextb"

		return input, nil
	}

	subscribers := map[string]xixo.CallbackMap{"root": callback}
	driver := xixo.NewDriver(reader, &writer, subscribers)

	// Stream the XML using the driver, assert the expected output, and check if the callback was called.
	err := driver.Stream()
	assert.Nil(t, err)

	assert.True(t, called)

	expected := "<root><element1>innerTextb</element1></root>"
	assert.Equal(t, expected, writer.String())
}

func TestFuncDriverEditEmptyElement(t *testing.T) {
	t.Parallel()

	// Create a reader with an XML string, an empty writer, a callback function, and a driver.
	reader := bytes.NewBufferString("<root><element1 nil=\"true\"/></root>")
	writer := bytes.Buffer{}
	called := false
	callback := func(input map[string]string) (map[string]string, error) {
		called = true

		return input, nil
	}

	subscribers := map[string]xixo.CallbackMap{"root": callback}
	driver := xixo.NewDriver(reader, &writer, subscribers)

	// Stream the XML using the driver, assert the expected output, and check if the callback was called.
	err := driver.Stream()
	assert.Nil(t, err)

	assert.True(t, called)

	expected := "<root><element1 nil=\"true\"/></root>"
	assert.Equal(t, expected, writer.String())
}

func TestFuncDriverEdit2subscribers(t *testing.T) {
	t.Parallel()

	// Create a reader with an XML string, an empty writer, a callback function, and a driver.
	reader := bytes.NewBufferString(
		`<root>
			<root1><element1>innerTexta1</element1></root1>
			<root2><element2>innerTexta2</element2></root2>
		</root>`,
	)
	writer := bytes.Buffer{}
	called1, called2 := false, false

	subscribers := map[string]xixo.CallbackMap{
		"root1": func(input map[string]string) (map[string]string, error) {
			called1 = true
			assert.Equal(t, input["element1"], "innerTexta1")
			input["element1"] = "innerTextb1"

			return input, nil
		},
		"root2": func(input map[string]string) (map[string]string, error) {
			called2 = true
			assert.Equal(t, "innerTexta2", input["element2"])
			input["element2"] = "innerTextb2"

			return input, nil
		},
	}

	driver := xixo.NewDriver(reader, &writer, subscribers)

	// Stream the XML using the driver, assert the expected output, and check if the callback was called.
	err := driver.Stream()
	assert.Nil(t, err)

	assert.True(t, called1)
	assert.True(t, called2)

	expected := `<root>
			<root1><element1>innerTextb1</element1></root1>
			<root2><element2>innerTextb2</element2></root2>
		</root>`
	assert.Equal(t, expected, writer.String())
}
