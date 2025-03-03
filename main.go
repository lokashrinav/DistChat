package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -lstdc++
#include "encryption.h"
#include <stdlib.h>
*/
import "C"
import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"unsafe"
)

var (
	c = make(map[net.Conn]bool)
	b = make(chan string)
	m = &sync.Mutex{}
	k = "secret"
)

func eM(s string) string {
	cs := C.CString(s)
	ck := C.CString(k)
	defer C.free(unsafe.Pointer(cs))
	defer C.free(unsafe.Pointer(ck))
	r := C.en(cs, ck)
	rs := C.GoString(r)
	C.fr(r)
	return rs
}

func dM(s string) string {
	cs := C.CString(s)
	ck := C.CString(k)
	defer C.free(unsafe.Pointer(cs))
	defer C.free(unsafe.Pointer(ck))
	r := C.de(cs, ck)
	rs := C.GoString(r)
	C.fr(r)
	return rs
}

func h(conn net.Conn) {
	defer conn.Close()
	m.Lock()
	c[conn] = true
	m.Unlock()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		s := scanner.Text()
		d := dM(s)
		b <- fmt.Sprintf("%s: %s", conn.RemoteAddr().String(), d)
	}
	m.Lock()
	delete(c, conn)
	m.Unlock()
}

func br() {
	for {
		s := <-b
		e := eM(s)
		m.Lock()
		for conn := range c {
			fmt.Fprintln(conn, e)
		}
		m.Unlock()
	}
}

func main() {
	ln, _ := net.Listen("tcp", ":9000")
	defer ln.Close()
	go br()
	for {
		conn, _ := ln.Accept()
		go h(conn)
	}
}
