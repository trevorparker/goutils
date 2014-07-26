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
	show_tabs             bool
}

const usage_message string = "usage: cat [OPTION ...] [FILE ...]"
const help_message string = `Concatenate and print FILE or STDIN to STDOUT.

  -b, --number-nonblank     number only non-blank lines
  -E, --show-ends           print $ at the end of each line
  -n, --number              number output lines, starting with 1
  -s, --squeeze-blank       print no more than one consecutive blank line
  -T, --show-tabs           print tab character as ^I
  -h, --help                print this help message and exit
`

const tab rune = 9
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
	w := bufio.NewWriterSize(os.Stdout, 512)

	line_number := 0
	newline_next := true

	prev_rune := []rune{0, 0}

	buf := make([]byte, 512)

	for {
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
					fmt.Fprintf(w, "%6d\t", line_number)
					newline_next = false
				}
			}

			if args.show_tabs && this_rune == tab {
				w.Write([]byte("^I"))
				continue
			} else if this_rune == newline {
				newline_next = true
				if args.show_line_endings {
					w.Write([]byte("$"))
				}
			}

			w.Write([]byte(buf[i : i+1]))
		}

		w.Flush()
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
			if os.Args[i] == "-T" || os.Args[i] == "--show-tabs" {
				args.show_tabs = true
				continue
			}
			if os.Args[i] == "--" {
				reached_files = true
				continue
			}
			if os.Args[i] == "-" {
				reached_files = true
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
			if args.file[i] == "-" {
				cat(nil, args)
			} else {
				file, err := os.Open(args.file[i])
				if err != nil {
					panic(err)
				}
				cat(file, args)
			}
		}
	}
}
