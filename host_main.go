package main

/*
#include <sys/types.h>       // for __kernel_sa_family_t / sa_family_t
#include <sys/socket.h>      // for AF_VSOCK, SOCK_STREAM, struct sockaddr
#include <linux/vm_sockets.h> // for struct sockaddr_vm, VMADDR_CID_HOST
#include <unistd.h>           // for read/write/close
#include <stdint.h>           // for uint32_t
*/
import "C"
import (
    "fmt"
    "log"
    "unsafe"
)

func main() {
    const (
        enclaveCID  = C.uint32_t(16)   // must match the CID you used when launching
        enclavePort = C.uint32_t(5005) // must match the port your enclave is listening on
    )

    // 1) Create a VSOCK socket
    fd := C.socket(C.AF_VSOCK, C.SOCK_STREAM, 0)
    if fd < 0 {
        log.Fatalf("Host: socket() failed: %d", fd)
    }
    defer C.close(fd)

    // 2) Build the sockaddr_vm struct
    var addr C.struct_sockaddr_vm
    addr.svm_family = C.sa_family_t(C.AF_VSOCK)
    addr.svm_cid    = enclaveCID
    addr.svm_port   = enclavePort

    // 3) Connect to the enclave
    if rc := C.connect(
        fd,
        (*C.struct_sockaddr)(unsafe.Pointer(&addr)),
        C.socklen_t(unsafe.Sizeof(addr)),
    ); rc != 0 {
        log.Fatalf("Host: connect() failed: %d", rc)
    }

    // 4) Send a greeting
    msg := "Hello from the host!"
    if n := C.write(fd, unsafe.Pointer(&[]byte(msg)[0]), C.size_t(len(msg))); n < 0 {
        log.Fatalf("Host: write() failed: %d", n)
    }
    fmt.Printf("Host: sent: %q\n", msg)

    // 5) Read the response
    buf := make([]byte, 1024)
    n := C.read(fd, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
    if n < 0 {
        log.Fatalf("Host: read() failed: %d", n)
    }
    fmt.Printf("Host: received: %q\n", string(buf[:n]))
}