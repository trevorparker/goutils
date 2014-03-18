// cat -- concatenate and print files
// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2014 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the Modified BSD License (see LICENSE)

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

type arg struct {
	file         []string
	line_numbers bool
}

const usage_message string = "usage: cat [OPTION ...] [FILE ...]"
const help_message string = `Concatenate and print FILE or STDIN to STDOUT.

  -n, --number              number output lines, starting with 1
  -h, --help                print this help message and exit
`

const newline rune = 10

func usage(error string) {
	fmt.Fprintf(os.Stderr, "cat: %s\n%s\n", error, usage_message)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s\n%s", usage_message, help_message)
	os.Exit(0)
}

func cat(file io.Reader, args arg) {
	if file == nil {
		file = os.Stdin
	}
	r := bufio.NewReader(file)

	line_number := 0
	newline_next := true

	for {
		buf := make([]byte, 16)
		n, err := r.Read(buf)

		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}

		for i := 0; i < n; i++ {
			this_rune, _ := utf8.DecodeRune(buf[i : i+1])

			if args.line_numbers && newline_next == true {
				line_number++
				fmt.Printf("%6d\t", line_number)
				newline_next = false
			}

			if this_rune == newline {
				newline_next = true
			}

			os.Stdout.Write(buf[i : i+1])
		}
	}
}

func main() {
	args := arg{[]string{}, false}
	reached_files := false

	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			if os.Args[i] == "-n" || os.Args[i] == "--number" {
				args.line_numbers = true
				continue
			}
		}
		reached_files = true
		arg_v := os.Args[i]
		args.file = append(args.file, arg_v)
	}

	if len(args.file) == 0 {
		cat(nil, args)
	} else {
		for i := range args.file {
			file, err := os.Open(args.file[i])
			if err != nil {
				panic(err)
			}
			cat(file, args)
		}
	}
}
