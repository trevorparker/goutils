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

func wc(file io.Reader, args arg, size int64) int64 {
	// Default the count to the size passed in (i.e.: stat() on a file)
	c := size

	if file == nil {
		file = os.Stdin
	}

	// Count bytes coming in through STDIN
	if c == 0 && args.count_bytes {
		s := bufio.NewScanner(file)
		s.Split(bufio.ScanBytes)
		for s.Scan() {
			c++
		}
	}

	// Count lines or words on STDIN or in the file passed in
	if args.count_lines && !args.count_bytes {
		c = 0
		s := bufio.NewScanner(file)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			c++
		}
	} else if args.count_words && !args.count_bytes {
		c = 0
		s := bufio.NewScanner(file)
		s.Split(bufio.ScanWords)
		for s.Scan() {
			c++
		}
	}

	return c
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
				args.count_lines = true
				continue
			}
			if os.Args[i] == "-w" || os.Args[i] == "--words" {
				args.count_words = true
				continue
			}
		}
		arg_v := os.Args[i]
		reached_files = true
		args.file = append(args.file, arg_v)
	}

	if len(args.file) == 0 {
		count := wc(nil, args, 0)
		fmt.Fprintf(os.Stdout, "%d\n", count)
	} else {
		for i := range args.file {
			// Call stat() on the file to help byte count performance
			file, err := os.Open(args.file[i])
			stat, err := file.Stat()
			size := stat.Size()
			if err != nil {
				panic(err)
			}

			count := wc(file, args, size)
			fmt.Fprintf(os.Stdout, "%d %s\n", count, args.file[i])
		}
	}
}
