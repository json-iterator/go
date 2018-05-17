package jsoniter

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFrozenConfig_MarshalIndentErrors(t *testing.T) {
	var panicOccured interface{} = nil

	defer func() {
		panicOccured = recover()
	}()
	_, err := MarshalIndent(`{"foo":"bar"}`, "", "\n\n")

	assert.NoError(t, err)
	assert.NotNil(t, panicOccured, "MarshalIndent should panic due to invalid indent chars")

	panicOccured = nil
	_, err = MarshalIndent(`{"foo":"bar"}`, "", " \t ")
	panicOccured = recover()
	assert.NoError(t, err)
	assert.NotNil(t, panicOccured, "MarshalIndent should panic due to mixed spaces and tabs")
}
