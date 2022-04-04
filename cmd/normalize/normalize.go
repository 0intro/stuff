// normalize - normalize unicode strings
//
// Copyright 2013 David du Colombier.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var normForm = flag.String("f", "nfc", "normalization form (NFC, NFD, NFKC or NFKD)")

var forms = map[string]norm.Form{
	"nfc":  norm.NFC,
	"nfd":  norm.NFD,
	"nfkc": norm.NFKC,
	"nfkd": norm.NFKD,
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: normalize [ -f form ] [ file ... ]\n")
	os.Exit(2)
}

func sysfatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "normalize: %s\n", fmt.Sprintf(format, args...))
	os.Exit(2)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		normalize("<stdin>", os.Stdin)
	} else {
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				sysfatal(err.Error())
			}
			normalize(arg, f)
			f.Close()
		}
	}
}

func form(name string) norm.Form {
	if f, ok := forms[strings.ToLower(name)]; ok {
		return f
	}
	sysfatal("invalid normalization form")
	return 0
}

func normalize(name string, r io.Reader) {
	buf := make([]byte, 8192)
	w := form(*normForm).Writer(os.Stdout)
	defer w.Close()
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			sysfatal(err.Error())
		}
		if n == 0 {
			break
		}
		_, err = w.Write(buf[:n])
		if err != nil {
			sysfatal(err.Error())
		}
	}
}
