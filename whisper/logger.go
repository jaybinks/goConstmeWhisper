package whisper

import (
	"C"
	"fmt"
)

/*
	https://github.com/Const-me/Whisper/blob/843a2a6ca6ea47c5ac4889a281badfc808d0ea01/Whisper/API/loggerApi.h

*/

type eLogLevel uint8

const (
	LlError   eLogLevel = 0
	LlWarning           = 1
	LlInfo              = 2
	LlDebug             = 3
)

type eLogFlags uint8

const (
	LfUndocumented      eLogFlags = 0
	LfUseStandardError            = 1
	LfSkipFormatMessage           = 2
)

type sLoggerSetup struct {
	sink    uintptr   // pfnLoggerSink
	context uintptr   // void*
	level   eLogLevel // eLogLevel
	flags   eLogFlags // eLoggerFlags
}

func initDefaultLogger() sLoggerSetup {
	thislogger := sLoggerSetup{}

	thislogger.sink = 0
	thislogger.context = 0
	thislogger.level = LlDebug
	thislogger.flags = LfUseStandardError

	/*
		cb := fnLoggerSink
		thislogger.sink = syscall.NewCallback(cb)
	*/

	return thislogger
}

func fnLoggerSink(context uintptr, lvl eLogLevel, message *C.char) uintptr {

	strmessage := C.GoString(message)
	fmt.Printf("%d - %s\n", lvl, strmessage)

	return 0
}
