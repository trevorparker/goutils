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
	file                  []string
	nonblank_line_numbers bool
	line_numbers          bool
	show_line_endings     bool
	squeeze_blank         bool
}

const usage_message string = "usage: cat [OPTION ...] [FILE ...]"
const help_message string = `Concatenate and print FILE or STDIN to STDOUT.

  -b, --number-nonblank     number only non-blank lines
  -E, --show-ends           print $ at the end of each line
  -n, --number              number output lines, starting with 1
  -s, --squeeze-blank       print no more than one consecutive blank line
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

	prev_rune := []rune{0, 0}

	for {
		buf := make([]byte, 16)
		n, err := r.Read(buf)

		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}

		for i := 0; i < n; i++ {
			this_rune, _ := utf8.DecodeRune(buf[i : i+1])

			// Identify and squeeze blank lines
			if args.squeeze_blank && prev_rune[0] == newline {
				if this_rune == newline && prev_rune[0] == prev_rune[1] {
					continue
				}
			}

			prev_rune = []rune{this_rune, prev_rune[0]}

			if args.line_numbers && newline_next == true {
				if this_rune != newline || !args.nonblank_line_numbers {
					line_number++
					fmt.Printf("%6d\t", line_number)
					newline_next = false
				}
			}

			if this_rune == newline {
				newline_next = true
				if args.show_line_endings {
					os.Stdout.Write([]byte("$"))
				}
			}

			os.Stdout.Write(buf[i : i+1])
		}
	}
}

func main() {
	args := arg{}
	reached_files := false

	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			if os.Args[i] == "-b" || os.Args[i] == "--number-nonblank" {
				args.nonblank_line_numbers = true
				args.line_numbers = true
				continue
			}
			if os.Args[i] == "-E" || os.Args[i] == "--show-ends" {
				args.show_line_endings = true
				continue
			}
			if os.Args[i] == "-n" || os.Args[i] == "--number" {
				args.line_numbers = true
				continue
			}
			if os.Args[i] == "-s" || os.Args[i] == "--squeeze-blank" {
				args.squeeze_blank = true
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
