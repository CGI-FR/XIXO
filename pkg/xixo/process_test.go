package xixo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youen/xixo/pkg/xixo"
)

func TestProcessCallbackWriteParamToStdinAndReturnStdout(t *testing.T) {
	t.Parallel()

	process := xixo.NewProcess("tr 1 2")

	err := process.Start()
	assert.Nil(t, err)

	result, err := process.Callback()("element1")
	assert.Nil(t, err)

	err = process.Stop()
	assert.Nil(t, err)

	assert.Equal(t, "element2", result)
}

func TestProcessCallbackWriteParamToStdinAndReturnStdoutLimitTo2(t *testing.T) {
	t.Parallel()

	process := xixo.NewProcess("head -n 2 | tr 1 2")

	err := process.Start()
	assert.Nil(t, err)

	for i := 0; i < 2; i++ {
		result, err := process.Callback()("element1")
		assert.Nil(t, err)
		assert.Equal(t, "element2", result)
	}

	_, err = process.Callback()("element1")
	assert.NotNil(t, err)

	err = process.Stop()
	assert.Nil(t, err)
}

func TestProcessCallbackFailedToStart(t *testing.T) {
	t.Parallel()

	process := xixo.NewProcess("false")

	err := process.Start()
	assert.Nil(t, err)

	_, err = process.Callback()("element1")
	assert.NotNil(t, err)

	err = process.Stop()
	assert.Nil(t, err)
}
