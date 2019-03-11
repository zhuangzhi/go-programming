# CGO related

https://golang.org/cmd/cgo/

## 1. Simple Example

```go

// #include <stdio.h>
// #include <stdint.h>
// #include <stdlib.h>
// #include <string.h>
//
// char* print(void* data)  {
//    printf("%s\n", (char*)data);
//    char *out = (char*)malloc(strlen(data)+1)
//    strcpy(out, data);
//    return out;
// };
import "C"

// TestC testing 'C' used func.
func PrintC(data string) {
    in := C.CString(data) // It allocate a C memory, user should free it.
    out := C.print(in)
    fmt.Printf("ret:%v\n", C.GoString(out))
    C.free(unsafe.Pointer(in))
    C.free(unsafe.Pointer(out))
}


```

## 2. Swap data with C

### 2.1. Standard way

`C.CString()` function is used to convert Go string to C string:

```go
cstr := C.String(gostr)
```

Be aware that C.CString function allocate memory, but not released. User `MUST` release it.

```go
C.free(unsafe.Pointer(cstr))
```

`C.GoString` function is used to convert C string to Go string. User `MUST` release C generated memory as the C code can't release it automatically.

### 2.2. Fast way

Golang use a lot of CPU to invoke `C.CString` and `C.GoString`. We can use an unsafe way to just send the pointer of Go to C:

```go
// type cbuf struct {
//     void *data;
//     int length;
//     int cap;
//     int hold;
// }
// cbuf* new_cbuf() {
//     return (cbuf*)malloc(sizeof(cbuf));
// }
// void delete_cbuf(cbuf*buf) {
//     cbuf_free(buf);
//     free(buf);
// }
// void cbuf_free(cbuf* buf) {
//     if (buf->hold!=0) {
//         free(buf->data);
//     }
//     buf->data = 0;
//     buf->length = 0;
//     buf->cap = 0;
//     buf->hold = 0;
// }
// void cbuf_resolve(cbuf* buf, int length) {
//     if (buf->cap < length) {
//         cbuf_free(buf);
//     }
//     buf->data = malloc(length);
//     buf->length = length;
//     buf->cap = length;
//     buf->hold = 1;
// }
// void print(cbuf* data, cbuf *out)  {
//    // do sth with data.
//    const char* _const_out = "hello123";
//    const int out_len = strlen(_const_out)+1;
//    cbuf_resolve(out, out_len);
//    strcpy(out->data, _const_out);
// };
import "C"

func PrintC(data string) []byte {
    in, out:= C.new_cbuf(), C.new_cbuf()
    in.data = unsafe.Pointer(&data[0]);
    in.length = C.int(len(data))
    in.cap = C.int(len(data))
    outbuf := make([]byte, 1024)
    out.data = unsafe.Pointer(&outbuf[0])
    out.length = C.int(1024)
    out.cap = C.int(1024)
    C.print(in, out)
    defer func() {
        C.delete_cbuf(in)
        C.delete_cbuf(out)
    }()
    if int(out.length) > len(outbuf) {
        return C.GoString((C.char*)out.data)
    }
}
```

In this case we monitor the CPU usage, it don't have significant improvement. Why? Golang doing dynamic checks of memory exchange between Go and C.

### 2.3. Passing pointers

https://golang.org/cmd/cgo/ : Passing pointers

These rules are checked dynamically at runtime. The checking is controlled by the cgocheck setting of the GODEBUG environment variable. The default setting is `GODEBUG=cgocheck=1`, which implements reasonably cheap dynamic checks. These checks may be disabled entirely using `GODEBUG=cgocheck=0`. Complete checking of pointer handling, at some cost in run time, is available via `GODEBUG=cgocheck=2`.

When we set `GODEBUG=cgocheck=0` before compile go code, the cgo invoke become mush faster.

## 3. Thread module

Golang would create a thread pool to running C code. The CGO invoke `WOULD NOT` block the golang thread. If we set `runtime.GOMAXPROCS(1)`, and invoke a C blocking function, Golang go-routines would be running without blocking as the CGO create separate thread for C code running. If user use a lot of go-routine to invoke C functions, you can see there are a lot of thread is created for C code running (`ps -T -p {pid}`, `top -H`, `htop`). User should limit the go-routine number which invoke C api, if you want to limit the concurrency of C function invoke.

## 4. Dynamic Load libraries

Golang has a feature 'plugin' to dynamic load libraries in runtime.