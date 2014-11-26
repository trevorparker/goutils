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
)

type arg struct {
	file         []string
	one_per_line bool
}

const (
	usage_message string = "usage: ls [OPTION ...] [FILE ...]"
	help_message  string = `List files and directories, and information about them.

  -1             print one entry per line
  -h, --help     print this help message and exit
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

	if args.one_per_line {
		for _, e := range *entries {
			out.WriteString(fmt.Sprintf("%s\n", e.Name()))
		}
		fmt.Print(out.String())
	} else {
		longest_entry := 0
		for _, e := range *entries {
			length := len(e.Name())
			if length > longest_entry {
				longest_entry = length + 1
			}
		}

		columns := int(78 / longest_entry)
		formatted_string := fmt.Sprintf("%%-%ds", longest_entry)
		for i, e := range *entries {
			out.WriteString(fmt.Sprintf(formatted_string, e.Name()))
			if i%columns == columns-1 {
				out.WriteString("\n")
			}
		}
		fmt.Println(out.String())
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
