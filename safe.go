package utils

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"runtime"
	"utils/zaplogger"
)

type SafeGoroutineFunc func()

func SafeGoPanicHandler(err interface{}, logger *zaplogger.ZapLogger) {

	info := fmt.Sprintf("catch PANIC(%v), CALLSTACK list:\n", err)
	info = info + *RetrieveCallStack()
	if logger != nil {
		logger.Panic(info)
	} else {
		fmt.Fprint(os.Stderr, info)
	}
}

func SafeGoPanicHandlerWithExtraInfo(err interface{}, extraInfo *string, logger *zaplogger.ZapLogger) {

	extra := ""
	if extraInfo != nil {
		extra = *extraInfo
	}
	info := fmt.Sprintf("catch PANIC(%v), EXTRA info(%s), \nCALLSTACK list:\n", err, extra)
	info = info + *RetrieveCallStack()
	if logger != nil {
		logger.Panic(info)
	} else {
		fmt.Fprint(os.Stderr, info)
	}
}

func SafeGoPanicHttpHandler(err interface{}, logger *zaplogger.ZapLogger, request *string) {

	info := fmt.Sprintf("httpHandler catch PANIC(%v), CALLSTACK list:\n", err)
	info = info + *RetrieveCallStack()
	if logger != nil {
		logger.Panic(info, zap.String("_http_request", *request))
	} else {
		fmt.Fprint(os.Stderr, info)
	}
}

func SafeGoroutineStartWith(routine SafeGoroutineFunc, logger *zaplogger.ZapLogger) {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				SafeGoPanicHandler(err, logger)
			}
		}()
		routine()
	}()
}

func RetrieveCallStack() *string {

	info := ""
	for depth := 1; ; depth++ {
		add, fn, line, ok := runtime.Caller(depth)
		if ok == false {
			break
		}
		info = info + fmt.Sprintf("    > %08x %s:%d\n", add, fn, line)
	}

	return &info
}
