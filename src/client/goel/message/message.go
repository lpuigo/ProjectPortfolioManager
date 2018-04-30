package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

var (
	duration int    = 3000
	msgclass string = "message"
)

func SetDuration(msec int) {
	duration = msec
}

func messageString(vm *hvue.VM, msgtype, msg string, close bool) {
	vm.Call("$message", js.M{
		"showClose":   close,
		"message":     msg,
		"type":        msgtype,
		"duration":    duration,
		"customClass": msgclass,
	})
}

func InfoStr(vm *hvue.VM, msg string, close bool) {
	pdur := duration
	if close {
		duration = 0
	}
	messageString(vm, "info", msg, close)
	duration = pdur
}

func SuccesStr(vm *hvue.VM, msg string, close bool) {
	messageString(vm, "success", msg, close)
}

func WarningStr(vm *hvue.VM, msg string, close bool) {
	messageString(vm, "warning", msg, close)
}

func ErrorStr(vm *hvue.VM, msg string, close bool) {
	pdur := duration
	if close {
		duration = 0
	}
	messageString(vm, "error", msg, close)
	duration = pdur
}
