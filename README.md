# goConstmeWhisper

Go Bindings for Const-Me's High-performance GPGPU Whisper implementation

https://github.com/Const-me/Whisper/

This does NOT use CGO, so does not require GCC on Windows.

Status is Working, but not complete. (It works for my usecase)
Untested on anything other than windows, but rumors suggest it may work in wine ?! 

# Todo Items
## General
- Wrap whisper.go in a class
- check whisper.dll exists, protect from errors loading the dll
- lots of tidyup
- Cleanup syscalls to all be SyscallN

## Testing
- How about we have some?

## IMediaFoundation
- Properly implement stero in LoadAudio* calls


