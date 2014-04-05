// tlsproxy - trivial tls proxy
//
// Copyright 2013 David du Colombier.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

type Proxy struct {
	cert       tls.Certificate
	clientTls  bool
	serverTls  bool
	nextprotos []string
}

func copy(raddr io.WriteCloser, laddr io.Reader) {
	defer raddr.Close()
	io.Copy(raddr, laddr)
}

func (proxy Proxy) Listen(laddr string, raddr string) {
	var err error
	var local net.Listener

	if proxy.serverTls {
		config := &tls.Config{NextProtos: proxy.nextprotos, Certificates: []tls.Certificate{proxy.cert}}
		local, err = tls.Listen("tcp", laddr, config)
	} else {
		local, err = net.Listen("tcp", laddr)
	}

	if err != nil {
		log.Fatal("error: listen:", err)
	}

	for {
		conn, err := local.Accept()
		if err != nil {
			log.Println("error: accept:", err)
		}
		go proxy.Serve(conn, raddr)
	}
}

func (proxy Proxy) Serve(local net.Conn, raddr string) {
	var err error
	var remote net.Conn

	if proxy.clientTls {
		config := &tls.Config{NextProtos: proxy.nextprotos, InsecureSkipVerify: true}
		remote, err = tls.Dial("tcp", raddr, config)
	} else {
		remote, err = net.Dial("tcp", raddr)
	}

	if err != nil {
		log.Println("error: dial:", err)
		return
	}

	go copy(local, remote)
	go copy(remote, local)
}

func NewProxy(c bool, s bool, crt string, key string) *Proxy {
	proxy := &Proxy{clientTls: c, serverTls: s}

	if proxy.serverTls {
		if cert, err := tls.LoadX509KeyPair(crt, key); err == nil {
			proxy.cert = cert
		} else {
			log.Fatal("error: LoadX509KeyPair:", err)
		}
	}

	return proxy
}
