// echo -- print arguments to standard output
// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2013-2014 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the Modified BSD License (see LICENSE)

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const usage_message string = "usage: echo [OPTION ...] [STRING ...]"
const help_message string = `Print STRING arguments to STDOUT.
Backslash escape sequences in STRING are interpreted.

  -n,                       do not print a trailing newline character
  -h, --help                print this help message and exit
`

func usage(error string) {
	fmt.Fprintf(os.Stderr, "echo: %s\n%s\n", error, usage_message)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s\n%s", usage_message, help_message)
	os.Exit(0)
}

func main() {
	start := 1
	trailing := "\n"
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		help()
	}
	if len(os.Args) > 1 && os.Args[1] == "-n" {
		start = 2
		trailing = ""
	}
	var b bytes.Buffer
	skip_next := false
	arg_string := strings.Join(os.Args[start:], " ")
	for i, c := range arg_string {
		if skip_next == true {
			skip_next = false
			continue
		} else {
			skip_next = true
		}
		r, _ := utf8.DecodeRune([]byte("\\"))
		if c == r && i < len(arg_string) {
			switch d, _ := utf8.DecodeRune([]byte(arg_string[i+1 : i+2])); d {
			case 'a':
				c, _ = utf8.DecodeRune([]byte("\a"))
			case 'b':
				c, _ = utf8.DecodeRune([]byte("\b"))
			case 'c':
				os.Stdout.Write([]byte(b.String()))
				os.Exit(0)
			case 'e':
				c, _ = utf8.DecodeRune([]byte("\x1B"))
			case 'f':
				c, _ = utf8.DecodeRune([]byte("\f"))
			case 'n':
				c, _ = utf8.DecodeRune([]byte("\n"))
			case 'r':
				c, _ = utf8.DecodeRune([]byte("\r"))
			case 't':
				c, _ = utf8.DecodeRune([]byte("\t"))
			case 'v':
				c, _ = utf8.DecodeRune([]byte("\v"))
			default:
				skip_next = false
			}
		} else {
			skip_next = false
		}
		b.WriteRune(c)
	}
	os.Stdout.Write([]byte(b.String() + trailing))
}
