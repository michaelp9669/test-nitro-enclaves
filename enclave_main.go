
package main

/*
#include <sys/types.h>        // for sa_family_t
#include <sys/socket.h>       // for AF_VSOCK, SOCK_STREAM, struct sockaddr
#include <linux/vm_sockets.h> // for struct sockaddr_vm, VMADDR_CID_ANY
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
    const port = C.uint32_t(5005)

    // 1) Create a VSOCK socket
    fd := C.socket(C.AF_VSOCK, C.SOCK_STREAM, 0)
    if fd < 0 {
        log.Fatalf("Enclave: socket() failed: %d", fd)
    }
    defer C.close(fd)

    // 2) Bind to any CID on port 5005
    var addr C.struct_sockaddr_vm
    addr.svm_family = C.sa_family_t(C.AF_VSOCK)
    addr.svm_cid    = C.VMADDR_CID_ANY
    addr.svm_port   = port

    if rc := C.bind(
        fd,
        (*C.struct_sockaddr)(unsafe.Pointer(&addr)),
        C.socklen_t(unsafe.Sizeof(addr)),
    ); rc != 0 {
        log.Fatalf("Enclave: bind() failed: %d", rc)
    }

    // 3) Listen for one connection
    if rc := C.listen(fd, 1); rc != 0 {
        log.Fatalf("Enclave: listen() failed: %d", rc)
    }
    log.Printf("Enclave: listening on vsock port %d…", port)

    // 4) Accept & handle
    for {
        connFd := C.accept(fd, nil, nil)
        if connFd < 0 {
            log.Printf("Enclave: accept() error: %d", connFd)
            continue
        }
        handleConnection(connFd)
    }
}

func handleConnection(fd C.int) {
    defer C.close(fd)

    // Read up to 1024 bytes
    buf := make([]byte, 1024)
    n := C.read(fd, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
    if n <= 0 {
        log.Printf("Enclave: read() error: %d", n)
        return
    }
    msg := string(buf[:n])
    fmt.Printf("Enclave: received from host: %q\n", msg)

    // Send a response
    response := "Hello from the enclave!"
    if _, err := C.write(fd, unsafe.Pointer(&[]byte(response)[0]), C.size_t(len(response))); err != nil {
        log.Printf("Enclave: write() error: %v", err)
    }
}