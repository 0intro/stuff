// unicode - interpret unicode characters
//
// Derived from Plan 9's /sys/src/cmd/unicode.c
// http://plan9.bell-labs.com/sources/plan9/sys/src/cmd/unicode.c
//
// Copyright (C) 2003, Lucent Technologies Inc. and others. All Rights Reserved.
// Portions Copyright 2013 David du Colombier.  All Rights Reserved.
// Distributed under the terms of the Lucent Public License Version 1.02
// See http://plan9.bell-labs.com/plan9/license.html

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	hex = "0123456789abcdefABCDEF"
)

var (
	numout = flag.Bool("numout", false, "force numeric output")
	text   = flag.Bool("text", false, "convert from numbers to running text")
	bout   *bufio.Writer
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: unicode { [-t] hex hex ... | hexmin-hexmax ... | [-n] char ... }\n")
	os.Exit(2)
}

func sysfatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "unicode: %s\n", fmt.Sprintf(format, args...))
	os.Exit(2)
}

func main() {
	bout = bufio.NewWriter(os.Stdout)
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		usage()
	}
	if !*numout && strings.ContainsRune(args[0], '-') {
		if err := ranges(&args); err != nil {
			sysfatal(err.Error())
		}
	} else if *numout || !strings.ContainsAny(hex, args[0]) {
		if err := nums(&args); err != nil {
			sysfatal(err.Error())
		}
	} else {
		if err := chars(&args); err != nil {
			sysfatal(err.Error())
		}
	}
}

func ranges(args *[]string) error {
	var i int
	for _, q := range *args {
		if i = strings.Index(q, "-"); i < 0 {
			goto err
		}
		if !strings.ContainsAny(hex, q[:i]) {
			goto err
		}
		min, err := strconv.ParseUint(q[:i], 16, 32)
		if err != nil || min > utf8.MaxRune || q[i] != '-' {
			goto err
		}
		q = q[i+1:]
		if !strings.ContainsAny(hex, q) {
			goto err
		}
		max, err := strconv.ParseUint(q, 16, 32)
		if err != nil || max > utf8.MaxRune || max < min {
			goto err
		}
		for i := 1; min <= max; min++ {
			fmt.Fprintf(bout, "%.6x %c", min, min)
			if min == max || (i&7) == 0 {
				fmt.Fprint(bout, "\n")
			} else {
				fmt.Fprint(bout, "\t")
			}
			i++
		}
	}
	bout.Flush()
	return nil
err:
	return fmt.Errorf("bad range")
}

func nums(args *[]string) error {
	var w int
	utferr := make([]byte, utf8.UTFMax)
	r := utf8.RuneError
	rsz := utf8.EncodeRune(utferr, r)
	for _, q := range *args {
		for i := 0; i < len(q); i += w {
			r, w = utf8.DecodeRune([]byte(q[i:]))
			if r == utf8.RuneError {
				if len(q[i:]) != rsz || q[i:] != string(utferr) {
					return fmt.Errorf("invalid utf string")
				}
			}
			fmt.Fprintf(bout, "%.6x\n", r)
		}
	}
	bout.Flush()
	return nil
}

func chars(args *[]string) error {
	for _, q := range *args {
		if strings.ContainsAny(hex, q) == false {
			goto err
		}
		m, err := strconv.ParseUint(q, 16, 32)
		if err != nil || m < 0 || m > utf8.MaxRune {
			goto err
		}
		fmt.Fprintf(bout, "%c", m)
		if !*text {
			fmt.Fprint(bout, "\n")
		}
	}
	bout.Flush()
	return nil
err:
	return fmt.Errorf("bad unicode value")
}
