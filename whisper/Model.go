package whisper

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// External - Go version of the struct
type Model struct {
	cStruct *_IModel
}

// Internal - C Version of the structs
type _IModel struct {
	lpVtbl *IModelVtbl
}

// https://github.com/Const-me/Whisper/blob/master/Whisper/API/iContext.cl.h
type IModelVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	createContext    uintptr //( iContext** pp ) = 0;
	isMultilingual   uintptr //() = 0;
	getSpecialTokens uintptr //( SpecialTokens& rdi ) = 0;
	stringFromToken  uintptr //( whisper_token token ) = 0;
}

func NewModel(cstruct *_IModel) *Model {
	this := Model{}
	this.cStruct = cstruct
	return &this
}

func (this *Model) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.cStruct.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this.cStruct)),
		0,
		0)
	return int32(ret)
}

func (this *Model) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.cStruct.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this.cStruct)),
		0,
		0)
	return int32(ret)
}

func (this *Model) createContext() (*iContext, error) {
	var context *iContext

	ret, _, err := syscall.Syscall(
		this.cStruct.lpVtbl.createContext,
		2, // Why was this 1, rather than 2 ?? 1 seemed to work fine
		uintptr(unsafe.Pointer(this.cStruct)),
		uintptr(unsafe.Pointer(&context)),
		0)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("createContext failed: %w", err.Error())
	}

	if windows.Handle(ret) != windows.S_OK {
		return nil, fmt.Errorf("loadModel failed: %w", err)
	}

	return context, nil
}

func (this *Model) IsMultilingual() bool {
	ret, _, _ := syscall.SyscallN(
		this.cStruct.lpVtbl.isMultilingual,
		uintptr(unsafe.Pointer(this.cStruct)),
	)

	return bool(windows.Handle(ret) == windows.S_OK)
}
