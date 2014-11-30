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
)

type arg struct {
	file            []string
	almost_all      bool
	comma_separated bool
	one_per_line    bool
}

const (
	usage_message string = "usage: ls [OPTION ...] [FILE ...]"
	help_message  string = `List files and directories, and information about them.

  -A, --almost-all     include entries beginning with a dot, except
                       implied . and ..
  -m                   print a comma-separated list of entries
  -1                   print one entry per line
  -h, --help           print this help message and exit
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

	if args.one_per_line {
		for _, e := range filtered_entries {
			out.WriteString(fmt.Sprintf("%s\n", e.Name()))
		}
		fmt.Print(out.String())
	} else if args.comma_separated {
		for i, e := range filtered_entries {
			out.WriteString(e.Name())
			if i < len(filtered_entries)-1 {
				out.WriteString(", ")
			}
		}
		fmt.Println(out.String())
	} else {
		longest_entry := 1
		for _, e := range filtered_entries {
			length := len(e.Name())
			if length > longest_entry {
				longest_entry = length + 1
			}
		}

		columns := int(78 / longest_entry)
		formatted_string := fmt.Sprintf("%%-%ds", longest_entry)
		for i, e := range filtered_entries {
			out.WriteString(fmt.Sprintf(formatted_string, e.Name()))
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
		filtered_entries = append(filtered_entries, e)
	}

	return filtered_entries
}

func main() {
	args := arg{}
	reached_files := false

	for i := 1; i < len(os.Args); i++ {
		if reached_files == false {
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				help()
			}
			if os.Args[i] == "-A" {
				args.almost_all = true
				continue
			}
			if os.Args[i] == "-m" {
				args.comma_separated = true
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
