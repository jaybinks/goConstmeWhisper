//go:build windows
// +build windows

package whisper

import (
	"C"
	"errors"
	"syscall"
	"unsafe"

	// Using lxn/win because its COM functions expose raw HRESULTs
	"golang.org/x/sys/windows"
)
import "fmt"

/*
	eModelImplementation - TranscribeStructs.h

	// GPGPU implementation based on Direct3D 11.0 compute shaders
	GPU = 1,

	// A hybrid implementation which uses DirectCompute for encode, and decodes on CPU
	// Not implemented in the published builds of the DLL. To enable, change BUILD_HYBRID_VERSION macro to 1
	Hybrid = 2,

	// A reference implementation which uses the original GGML CPU-running code
	// Not implemented in the published builds of the DLL. To enable, change BUILD_BOTH_VERSIONS macro to 1
	Reference = 3,
*/

var (
	dll = syscall.NewLazyDLL("whisper.dll") // Todo wrap this in a class, check file exists, handle errors ... you know, just a few things.. AKA Stop being lazy

	setupLogger           = dll.NewProc("setupLogger")
	loadModel             = dll.NewProc("loadModel")
	initMediaFoundation   = dll.NewProc("initMediaFoundation")
	findLanguageKeyW      = dll.NewProc("findLanguageKeyW")
	findLanguageKeyA      = dll.NewProc("findLanguageKeyA")
	getSupportedLanguages = dll.NewProc("getSupportedLanguages")
)

// https://learn.microsoft.com/en-us/windows/win32/seccrypto/common-hresult-values
// https://pkg.go.dev/golang.org/x/sys/windows
const (
	E_INVALIDARG                      = 0x80070057
	ERROR_HV_CPUID_FEATURE_VALIDATION = 0xC0350038
)

func SetupLogger(level eLogLevel, flags eLogFlags, cb *any) (bool, error) {

	setup := sLoggerSetup{}
	setup.sink = 0
	setup.context = 0
	setup.level = level
	setup.flags = flags

	if cb != nil {
		setup.sink = syscall.NewCallback(cb)
	}

	res, _, err := setupLogger.Call(uintptr(unsafe.Pointer(&setup)))

	return windows.Handle(res) == windows.S_OK, err
}

func LoadWhisperModel(path string) (*Model, error) {
	var modelptr *_IModel

	whisperpath, _ := windows.UTF16PtrFromString(path)

	setup := ModelSetup().AsCType()
	obj, _, _ := loadModel.Call(uintptr(unsafe.Pointer(whisperpath)), uintptr(unsafe.Pointer(setup)), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&modelptr)))

	if windows.Handle(obj) != windows.S_OK {
		fmt.Printf("loadModel failed: %s\n", syscall.Errno(obj).Error())
		return nil, fmt.Errorf("loadModel failed: %s", syscall.Errno(obj))
	}

	if modelptr == nil {
		return nil, errors.New("loadModel did not return a Model")
	}

	if modelptr.lpVtbl == nil {
		return nil, errors.New("loadModel method table is nil")
	}

	model := NewModel(modelptr)

	return model, nil

}

func DoinitMediaFoundation() (*IMediaFoundation, error) {

	var mediafoundation *IMediaFoundation

	// initMediaFoundation( iMediaFoundation** pp );
	obj, _, _ := initMediaFoundation.Call(uintptr(unsafe.Pointer(&mediafoundation)))

	if windows.Handle(obj) != windows.S_OK {
		fmt.Printf("initMediaFoundation failed: %s\n", syscall.Errno(obj).Error())
		return nil, fmt.Errorf("initMediaFoundation failed: %s", syscall.Errno(obj))
	}

	if mediafoundation.lpVtbl == nil {
		return nil, errors.New("initMediaFoundation method table is nil")
	}

	return mediafoundation, nil
}
