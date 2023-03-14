package whisper

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://github.com/Const-me/Whisper/blob/843a2a6ca6ea47c5ac4889a281badfc808d0ea01/Whisper/API/iMediaFoundation.h

type iMediaFoundation struct {
	lpVtbl *iMediaFoundationVtbl
}

type iMediaFoundationVtbl struct {
	QueryInterface     uintptr
	AddRef             uintptr
	Release            uintptr
	loadAudioFile      uintptr //( LPCTSTR path, bool stereo, iAudioBuffer** pp ) const;
	openAudioFile      uintptr // ( LPCTSTR path, bool stereo, iAudioReader** pp );
	listCaptureDevices uintptr // ( pfnFoundCaptureDevices pfn, void* pv );
	openCaptureDevice  uintptr // ( LPCTSTR endpoint, const sCaptureParams& captureParams, iAudioCapture** pp );
}

func (this *iMediaFoundation) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *iMediaFoundation) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *iMediaFoundation) LoadAudioFile(file string, stereo bool) (*iAudioBuffer, error) {

	var buffer *iAudioBuffer

	UTFFileName, _ := windows.UTF16PtrFromString(file)

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.loadAudioFile,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(UTFFileName)),
		uintptr(1),
		uintptr(unsafe.Pointer(&buffer)))

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("loadAudioFile failed: %s\n", syscall.Errno(ret).Error())
		return nil, syscall.Errno(ret)
	}

	return buffer, nil
}

// ************************************************************

type iAudioBuffer struct {
	lpVtbl *iAudioBufferVtbl
}

type iAudioBufferVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	countSamples   uintptr // returns uint32_t
	getPcmMono     uintptr // returns float*
	getPcmStereo   uintptr // returns float*
	getTime        uintptr // ( int64_t& rdi )
}

func (this *iAudioBuffer) AddRef() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioBuffer) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.Release,
		uintptr(unsafe.Pointer(this)),
	)
	return int32(ret)
}

func (this *iAudioBuffer) CountSamples() (uint32, error) {

	ret, _, err := syscall.SyscallN(
		this.lpVtbl.countSamples,
		uintptr(unsafe.Pointer(this)),
	)

	if err != 0 {
		return 0, errors.New(err.Error())
	}

	return uint32(ret), nil
}

// ************************************************************

type iAudioReader struct {
	lpVtbl *iAudioReaderVtbl
}

type iAudioReaderVtbl struct {
	QueryInterface  uintptr
	AddRef          uintptr
	Release         uintptr
	getDuration     uintptr // ( int64_t& rdi )
	getReader       uintptr // ( IMFSourceReader** pp )
	requestedStereo uintptr // ()
}

// ************************************************************

type iAudioCapture struct {
	lpVtbl *iAudioCaptureVtbl
}

type iAudioCaptureVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	getReader      uintptr // ( IMFSourceReader** pp )
	getParams      uintptr // returns sCaptureParams&
}
