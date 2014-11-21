// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2013-2014 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the Modified BSD License (see LICENSE)

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type arg struct {
	count   int
	bytes   int
	quiet   bool
	verbose bool
	file    []string
}

const (
	usage_message string = "usage: head [OPTION ...] [FILE ...]"
	help_message  string = `Print the front matter of FILE or STDIN.
A header describing the file name is prefixed when multiple files are passed
in. When no FILE is provided, read from STDIN.

  -c, --bytes=N             print the first N bytes of FILE or STDIN
  -n, --lines=N             print the first N lines of FILE or STDIN;
                                default 10
  -q, --quiet, --silent     don't print file name headers
  -v, --verbose             always print file name headers
  -h, --help                print this help message and exit
`
)

func usage(error string) {
	fmt.Fprintf(os.Stderr, "head: %s\n%s\n", error, usage_message)
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

func head(file io.Reader, args arg) {
	if file == nil {
		file = os.Stdin
	}

	r := bufio.NewReader(file)
	if args.bytes > 0 {
		// Create a buffer to fill with the number of bytes
		// requested
		var buffer bytes.Buffer
		for b := 0; b < args.bytes; b++ {
			c, err := r.ReadByte()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			buffer.WriteByte(c)
		}
		os.Stdout.Write([]byte(buffer.String()))
	} else {
		// Write out each line until we reach the number of
		// lines requested
		for l := 0; l < args.count; l++ {
			l, err := r.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			os.Stdout.Write([]byte(string(l)))
		}
	}
}

func main() {
	args := arg{10, 0, false, false, []string{}}
	reached_files := false
	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			var err error
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			arg_v := parse_args(os.Args, &i, "-n", "--count")
			if arg_v != "" {
				args.count, err = strconv.Atoi(arg_v)
				if err != nil {
					usage("illegal option " + os.Args[i])
				}
				continue
			}
			arg_v = parse_args(os.Args, &i, "-c", "--bytes")
			if arg_v != "" {
				args.bytes, err = strconv.Atoi(arg_v)
				if err != nil {
					usage("illegal option " + os.Args[i])
				}
				continue
			}
			if os.Args[i] == "-q" || os.Args[i] == "--quiet" || os.Args[i] == "--silent" {
				args.quiet = true
				continue
			}
			if os.Args[i] == "-v" || os.Args[i] == "--verbose" {
				args.verbose = true
				continue
			}
			arg_v = parse_args(os.Args, &i, "-", "-")
			if arg_v != "" {
				args.count, err = strconv.Atoi(arg_v)
				if err != nil {
					usage("illegal option " + os.Args[i])
				}
				continue
			}
		}
		arg_v := os.Args[i]
		reached_files = true
		args.file = append(args.file, arg_v)
	}

	if len(args.file) == 0 {
		head(nil, args)
	} else {
		for i := range args.file {
			// Print headers for the filenames if we are handling
			// multiple files
			if len(args.file) > 1 && !args.quiet || args.verbose {
				if i > 0 {
					fmt.Fprintf(os.Stdout, "\n==> %s <==\n", args.file[i])
				} else {
					fmt.Fprintf(os.Stdout, "==> %s <==\n", args.file[i])
				}
			}
			file, err := os.Open(args.file[i])
			if err != nil {
				panic(err)
			}
			head(file, args)
		}
	}
}
