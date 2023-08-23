package article

// #cgo CFLAGS: -I${SRCDIR}/../../include
// #cgo LDFLAGS: -static  lib/libsharpdown.a lib/libbootstrapperdll.a lib/libRuntime.WorkstationGC.a lib/libSystem.Native.a /usr/lib/libstdc++.a /usr/lib/libm.a -Wl,--require-defined,NativeAOT_StaticInitialization -Wl,--defsym,__xmknod=mknod
// #include <stdlib.h>
// #include "sharpdown.h"
import "C"
import "unsafe"

func Markdown2Html(markdown string) string {

	cMarkdown := C.CString(markdown)
	defer C.free(unsafe.Pointer(cMarkdown))

	cHtml := C.sharpdown_tohtml(cMarkdown)
	defer C.sharpdown_free(cHtml)

	return C.GoString(cHtml)
}
