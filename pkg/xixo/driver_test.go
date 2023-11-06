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

	expected := "<root>\n  <element1>innerTextb</element1>\n</root>"
	assert.Equal(t, expected, writer.String())
}
