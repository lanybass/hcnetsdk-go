package playctrl

import (
	"fmt"
	"syscall"
	"unsafe"
)

var PlayCtrlDLL = syscall.MustLoadDLL("PlayCtrl.dll")

func PlayM4_Play(port int32, hWnd unsafe.Pointer) {
	proc := PlayCtrlDLL.MustFindProc("PlayM4_Play")
	r, _, _ := proc.Call(uintptr(port), uintptr(hWnd))
	fmt.Printf("%v\n", r)
}
