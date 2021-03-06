// ls -- list files and directories
// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2014 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the Modified BSD License (see LICENSE)

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type arg struct {
	file            []string
	almost_all      bool
	ignore_backups  bool
	comma_separated bool
	quote_name      bool
	one_per_line    bool
}

const (
	usage_message string = "usage: ls [OPTION ...] [FILE ...]"
	help_message  string = `List files and directories, and information about them.

  -A, --almost-all      include entries beginning with a dot, except
                        implied . and ..
  -B, --ignore-backups  do not list entries ending with ~
  -m                    print a comma-separated list of entries
  -Q, --quote-name      print each entry surrounded by double quotes
  -1                    print one entry per line
  -h, --help            print this help message and exit
`
)

func usage(error string) {
	fmt.Fprintf(os.Stderr, "ls: %s\n%s\n", error, usage_message)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s\n%s", usage_message, help_message)
	os.Exit(0)
}

func ls(file string, args arg) {
	entries := make([]os.FileInfo, 0)

	// Determine if this is a file or directory, then call out
	// to ReadDir if it's a directory. Otherwise, we're can just
	// pass the file info on.
	fi, err := os.Stat(file)
	if err != nil {
		panic(err)
	} else if fi.IsDir() {
		e, err := ioutil.ReadDir(file)
		if err != nil {
			panic(err)
		}
		entries = e
	} else {
		entries = append(entries, fi)
	}
	printEntries(&entries, &args)
}

func printEntries(entries *[]os.FileInfo, args *arg) {
	var out bytes.Buffer

	filtered_entries := filterEntries(entries, args)

	// Determine the terminal width, useful for column and line
	// wrapping calculations.
	terminal_width, _, err := getTerminalSize()
	if err != nil {
		terminal_width = 78
	}

	if args.one_per_line {
		for _, e := range filtered_entries {
			name := e.Name()
			if args.quote_name {
				name = fmt.Sprintf("\"%s\"", e.Name())
			}
			out.WriteString(fmt.Sprintf("%s\n", name))
		}
		fmt.Print(out.String())
	} else if args.comma_separated {
		for i, e := range filtered_entries {
			var scratch bytes.Buffer
			name := e.Name()
			if args.quote_name {
				name = fmt.Sprintf("\"%s\"", e.Name())
			}
			scratch.WriteString(name)
			if i < len(filtered_entries)-1 {
				scratch.WriteString(", ")
			}

			// Finish out this line if we're going to hit the
			// terminal width. The next entry will wrap to the
			// next line.
			if out.Len()+scratch.Len() >= terminal_width {
				fmt.Println(out.String())
				out.Reset()
			}
			out.WriteString(scratch.String())
		}
		fmt.Println(out.String())
	} else {
		longest_entry := 1
		for _, e := range filtered_entries {
			name := e.Name()
			if args.quote_name {
				name = fmt.Sprintf("\"%s\"", e.Name())
			}
			length := len(name)
			if length > longest_entry {
				longest_entry = length + 1
			}
		}

		columns := int(terminal_width / longest_entry)

		formatted_string := fmt.Sprintf("%%-%ds", longest_entry)
		for i, e := range filtered_entries {
			name := e.Name()
			if args.quote_name {
				name = fmt.Sprintf("\"%s\"", e.Name())
			}
			out.WriteString(fmt.Sprintf(formatted_string, name))
			if i%columns == columns-1 {
				out.WriteString("\n")
			}
		}
		fmt.Println(out.String())
	}
}

func filterEntries(entries *[]os.FileInfo, args *arg) []os.FileInfo {
	filtered_entries := make([]os.FileInfo, 0)
	for _, e := range *entries {
		if !args.almost_all && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if args.ignore_backups && strings.HasSuffix(e.Name(), "~") {
			continue
		}
		filtered_entries = append(filtered_entries, e)
	}

	return filtered_entries
}

// This bit thanks in part to:
// - https://code.google.com/p/go/source/browse/ssh/terminal/util.go?repo=crypto#75 and
// - http://stackoverflow.com/questions/16569433/get-terminal-size-in-go
func getTerminalSize() (width, height int, err error) {
	var dimensions [4]uint16

	ret, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)))

	if ret != 0 {
		return -1, -1, err
	}

	return int(dimensions[1]), int(dimensions[0]), nil
}

func main() {
	args := arg{}
	reached_files := false

	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			if os.Args[i] == "-A" || os.Args[i] == "--almost-all" {
				args.almost_all = true
				continue
			}
			if os.Args[i] == "-B" || os.Args[i] == "--ignore-backups" {
				args.ignore_backups = true
				continue
			}
			if os.Args[i] == "-m" {
				args.comma_separated = true
				continue
			}
			if os.Args[i] == "-Q" || os.Args[i] == "--quote-name" {
				args.quote_name = true
				continue
			}
			if os.Args[i] == "-1" {
				args.one_per_line = true
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
		ls("./", args)
	} else {
		for i := range args.file {
			ls(args.file[i], args)
		}
	}
}
