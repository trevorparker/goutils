// head -- print the front matter of file(s) or standard input
// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2013 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the modified BSD license (see LICENSE)

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
	count int
	bytes int
	file  []string
}

const usage_message string = "usage: head [OPTION ...] [FILE ...]"
const help_message string = `Print the front matter of FILE or STDIN.
A header describing the file name is prefixed when multiple files are passed
in. When no FILE is provided, read from STDIN.

  -c, --bytes=N             print the first N bytes of FILE or STDIN
  -n, --lines=N             print the first N lines of FILE or STDIN;
                                default 10
  -h, --help                print this help message and exit
`

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
		print(buffer.String())
	} else {
		for l := 0; l < args.count; l++ {
			l, err := r.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			print(string(l))
		}
	}
}

func main() {
	args := arg{10, 0, []string{}}
	reached_files := false
	for i := 0; i < len(os.Args); i++ {
		if i == 0 {
			continue
		}
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
			if len(args.file) > 1 {
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
