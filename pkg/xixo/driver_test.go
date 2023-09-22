package xixo_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

// TestMapCallback should convert.
func TestDriverSimpleEdit(t *testing.T) {
	t.Parallel()

	reader := bytes.NewBufferString("<root><element1>innerTexta</element1></root>")
	writer := bytes.Buffer{}
	subscribers := map[string]string{"root": "tr a b"}
	driver := xixo.NewShellDriver(reader, &writer, subscribers)

	err := driver.Stream()
	assert.Nil(t, err)

	expected := "<root>\n  <element1>innerTextb</element1>\n</root>"
	assert.Equal(t, expected, writer.String())
}

func TestFuncDriverEdit(t *testing.T) {
	t.Parallel()

	reader := bytes.NewBufferString("<root><element1>innerTexta</element1></root>")
	writer := bytes.Buffer{}
	called := false
	callback := func(input map[string]string) (map[string]string, error) {
		called = true
		input["element1"] = "innerTextb"

		return input, nil
	}

	subscribers := map[string]xixo.CallbackMap{"root": callback}
	driver := xixo.NewFuncDriver(reader, &writer, subscribers)

	err := driver.Stream()
	assert.Nil(t, err)

	assert.True(t, called)

	expected := "<root>\n  <element1>innerTextb</element1>\n</root>"
	assert.Equal(t, expected, writer.String())
}
