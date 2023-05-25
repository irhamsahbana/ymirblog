// Package shared is library function for whole system.
// # This manifest was generated by ymir. DO NOT EDIT.
//go:build !windows
// +build !windows

package shared

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}