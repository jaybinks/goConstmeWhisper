package whisper

import (
	"C"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type eTokenFlags uint32

const (
	tfNone    eTokenFlags = 0
	tfSpecial             = 1
)

type sTranscribeLength struct {
	countSegments uint32
	countTokens   uint32
}

type sTimeSpan struct {

	// The value is expressed in 100-nanoseconds ticks: compatible with System.Timespan, FILETIME, and many other things
	ticks uint64

	/*
		operator sTimeSpanFields() const
		{
			return sTimeSpanFields{ ticks };
		}
		void operator=( uint64_t tt )
		{
			ticks = tt;
		}
		void operator=( int64_t tt )
		{
			assert( tt >= 0 );
			ticks = (uint64_t)tt;
		} */
}

type sTimeInterval struct {
	begin sTimeSpan
	end   sTimeSpan
}

type sSegment struct {
	// Segment text, null-terminated, and probably UTF-8 encoded
	text *C.char

	// Start and end times of the segment
	time sTimeInterval

	// These two integers define the slice of the tokens in this segment, in the array returned by iTranscribeResult.getTokens method
	firstToken  uint32
	countTokens uint32
}

type sSegmentArray []sSegment

type sToken struct {
	// Token text, null-terminated, and usually UTF-8 encoded.
	// I think for Chinese language the models sometimes outputs invalid UTF8 strings here, Unicode code points can be split between adjacent tokens in the same segment
	// More info: https://github.com/ggerganov/whisper.cpp/issues/399
	text *C.char

	// Start and end times of the token
	time sTimeInterval
	// Probability of the token
	probability float32

	// Probability of the timestamp token
	probabilityTimestamp float32

	// Sum of probabilities of all timestamp tokens
	ptsum float32

	// Voice length of the token
	vlen float32

	// Token id
	id int32

	flags eTokenFlags
}

type sTokenArray []sToken

type iTranscribeResultVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	getSize     uintptr // ( sTranscribeLength& rdi ) HRESULT
	getSegments uintptr // () getTokens
	getTokens   uintptr // () getToken*
}

type iTranscribeResult struct {
	lpVtbl *iTranscribeResultVtbl
}

func (this *iTranscribeResult) AddRef() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *iTranscribeResult) Release() int32 {
	ret, _, _ := syscall.Syscall(
		this.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(this)),
		0,
		0)
	return int32(ret)
}

func (this *iTranscribeResult) GetSize() (*sTranscribeLength, error) {

	var result sTranscribeLength

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getSize,
		uintptr(unsafe.Pointer(this)),
		uintptr(unsafe.Pointer(&result)),
	)

	if windows.Handle(ret) != windows.S_OK {
		fmt.Printf("iTranscribeResult.GetSize failed: %s\n", syscall.Errno(ret).Error())
		return nil, errors.New(syscall.Errno(ret).Error())
	}

	return &result, nil

}

func (this *iTranscribeResult) GetSegments(len uint32) []sSegment {

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getSegments,
		uintptr(unsafe.Pointer(this)),
	)

	data := unsafe.Slice((*sSegment)(unsafe.Pointer(ret)), len)

	return data
}

func (this *iTranscribeResult) GetTokens(len uint32) []sToken {

	ret, _, _ := syscall.SyscallN(
		this.lpVtbl.getTokens,
		uintptr(unsafe.Pointer(this)),
	)

	return unsafe.Slice((*sToken)(unsafe.Pointer(ret)), len)
}
