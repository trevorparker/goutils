// wc -- print byte, line, or word counts for files
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
	"strings"
)

type arg struct {
	count_bytes bool
	count_lines bool
	count_words bool
	file        []string
}

const (
	usage_message string = "usage: wc [OPTION ...] [FILE ...]"
	help_message  string = `Count bytes, lines, or words for FILE or STDIN to STDOUT.

  -c, --bytes     count bytes
  -l, --lines     count newlines
  -w, --words     count words
  -h, --help      print this help message and exit
`
)

func usage(error string) {
	fmt.Fprintf(os.Stderr, "wc: %s\n%s\n", error, usage_message)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s\n%s", usage_message, help_message)
	os.Exit(0)
}

func parse_args(args []string, i *int, s string, l string) (arg_v string) {
	if strings.HasPrefix(args[*i], s) || strings.HasPrefix(args[*i], l) {
		arg_v := strings.Trim(args[*i], s)
		if len(arg_v) == 0 && len(args)-1 > *i {
			*i++
			arg_v = args[*i]
		} else if len(arg_v) == 0 {
			usage("option requires value -- " + args[*i])
		}
		return arg_v
	}
	return ""
}

func wc(file io.Reader, args arg) int {
	if file == nil {
		file = os.Stdin
	}

	count := 0
	complete := false

	r := bufio.NewReader(file)
	for complete == false {
		b, err := r.ReadByte()
		if err != nil {
			complete = true
		} else if args.count_bytes {
			if b != 0 {
				count++
			}
		}
	}

	return count
}

func main() {
	args := arg{false, false, false, []string{}}
	reached_files := false
	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			if os.Args[i] == "-c" || os.Args[i] == "--bytes" {
				args.count_bytes = true
				continue
			}
			if os.Args[i] == "-l" || os.Args[i] == "--lines" {
				usage("not implemented")
				args.count_lines = true
				continue
			}
			if os.Args[i] == "-w" || os.Args[i] == "--words" {
				usage("not implemented")
				args.count_words = true
				continue
			}
		}
		arg_v := os.Args[i]
		reached_files = true
		args.file = append(args.file, arg_v)
	}
	if len(args.file) == 0 {
		count := wc(nil, args)
		fmt.Fprintf(os.Stdout, "%d\n", count)
	} else {
		for i := range args.file {
			file, err := os.Open(args.file[i])
			if err != nil {
				panic(err)
			}
			count := wc(file, args)
			fmt.Fprintf(os.Stdout, "%d %s\n", count, args.file[i])
		}
	}
}
