package main

import (
	"testing"

	"github.com/gookit/gcli/v3/progress"
	"github.com/stretchr/testify/assert"
)

func Test_CollectAll(t *testing.T) {
	p := progress.New(100)
	logger := NewLogger()
	p.Start()
	results := CollectAll(p, logger)
	p.Finish()
	assert.IsType(t, []*VendorDef{}, results)
	assert.True(t, len(results) > 5)
}
