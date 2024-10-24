//go:build darwin
package tunnel

import (
	"syscall"
	"fmt"
	"net"
)

func preAccept(conn *net.TCPListener) error {
    var readFDs syscall.FdSet
    var fd uintptr

	rawConn, err := conn.SyscallConn()
    if err != nil {
    	return fmt.Errorf("Failed to get raw connection : %s", err)
    }

    err = rawConn.Control(func(listenerFD uintptr) {
        fd = listenerFD
        FD_SET(int(listenerFD), &readFDs)
    })
    if err != nil {
        return fmt.Errorf("Failed to control connection: %w", err)
    }
    // Blocking until some connection ready to accept
    // Need to check timeout with loop ?
    err = syscall.Select(int(fd)+1, &readFDs, nil, nil, nil)
    if err != nil {
    	return fmt.Errorf("Select failed: %w", err)
    }
    return nil
}

func FD_SET(p int, set *syscall.FdSet) {
    set.Bits[p/64] |= 1 << (uint(p) % 64)
}

func FD_ISSET(p int, set *syscall.FdSet) bool {
    return (set.Bits[p/64] & (1 << (uint(p) % 64))) != 0
}