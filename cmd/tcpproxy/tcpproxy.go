// tcpproxy - trivial tcp proxy
//
// Copyright 2013 David du Colombier.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: tcpproxy local:port remote:port\n")
	os.Exit(2)
}

func sysfatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "tcpproxy: %s\n", fmt.Sprintf(format, args...))
	os.Exit(2)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	laddr := os.Args[1]
	raddr := os.Args[2]
	local, err := net.Listen("tcp", laddr)
	if local == nil {
		sysfatal("listen: %v", err)
	}
	for {
		conn, err := local.Accept()
		if conn == nil {
			sysfatal("accept: %v", err)
		}
		go forward(conn, raddr)
	}
}

func forward(local net.Conn, raddr string) {
	remote, err := net.Dial("tcp", raddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial: %v\n", err)
		return
	}
	go copy(local, remote)
	go copy(remote, local)
}

func copy(dst io.WriteCloser, src io.Reader) {
	defer dst.Close()
	io.Copy(dst, src)
}
