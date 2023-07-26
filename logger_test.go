package main

import (
	"io"
	"testing"
)

func Test_Logger(t *testing.T) {
	log := NewLogger()
	log.Success("test success")
	log.Success("test success formatting: %s", "this should be bold")
	log.Info("test info")
	log.Info("test info formatting: %s", "this should be bold")
	log.Warn("test warn")
	log.Warn("test warn formatting: %s", "this should be bold")
	log.Error("test error")
	log.Error("test error formatting: %s", "this should be bold")
	log.Err(io.EOF)
	log.Err(io.EOF, "this should bold")
	log.Err(io.EOF, "this should be %s", "formatted")
}
