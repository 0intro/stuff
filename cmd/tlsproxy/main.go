// tlsproxy - trivial tls proxy
//
// Copyright 2013 David du Colombier.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	cFlag = flag.Bool("c", false, "enable TLS on client")
	sFlag = flag.Bool("s", false, "enable TLS on server")

	nextprotos = flag.String("p", "http/1.1", "TLS Next Protocol Negotiation (http/1.1, spdy/3, etc.)")

	crt = flag.String("crt", "server.crt", "X.509 certificate")
	key = flag.String("key", "server.key", "X.509 private key")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tlsproxy -c [ -p nextprotos ] local:port remote:port")
	fmt.Fprintln(os.Stderr, "usage: tlsproxy -s [ -p nextprotos ] [ -crt crt.pem ] [ -key key.pem ] local:port remote:port")
	fmt.Fprintln(os.Stderr, "usage: tlsproxy -c -s [ -p nextprotos ] [ -crt crt.pem ] [ -key key.pem ] local:port remote:port")
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if flag.NArg() != 2 {
		usage()
	}

	if !*cFlag && !*sFlag {
		usage()
	}

	proxy := NewProxy(*cFlag, *sFlag, *crt, *key)
	proxy.nextprotos = strings.Split(*nextprotos, " ")
	proxy.Listen(args[0], args[1])
}
