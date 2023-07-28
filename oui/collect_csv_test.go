package oui_test

import (
	"testing"

	"github.com/gookit/gcli/v3/progress"
	"github.com/stretchr/testify/assert"
	"github.com/thatmattlove/oui/internal/logger"
	"github.com/thatmattlove/oui/oui"
)

func Test_CollectAll(t *testing.T) {
	p := progress.New(100)
	logger := logger.New()
	p.Start()
	results := oui.CollectAll(p, logger)
	p.Finish()
	assert.IsType(t, []*oui.VendorDef{}, results)
	assert.True(t, len(results) > 5)
}
