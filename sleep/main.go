// sleep -- suspend execution for a specified time
// Part of goutils (https://github.com/trevorparker/goutils)
//
// Copyright (c) 2014 Trevor Parker <trevor@trevorparker.com>
// All rights reserved
//
// Distributed under the terms of the Modified BSD License (see LICENSE)

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	usage_message string = "usage: sleep [OPTION ...] NUMBER[SUFFIX] ..."
	help_message  string = `Suspend execution for a a specified time.
Execution sleeps for a NUMBER of seconds. If multiple NUMBER arguments are
provided, execution will sleep for the sum of their durations. NUMBER may be
an integer or floating point number.

If SUFFIX is specified, execution will be suspended for a NUMBER of:
's': seconds; 'm': minutes; 'h': hours.

  -h, --help                print this help message and exit
`
)

func usage(error string) {
	fmt.Fprintf(os.Stderr, "sleep: %s\n%s\n", error, usage_message)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s\n%s", usage_message, help_message)
	os.Exit(0)
}

func sleep(duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		usage(err.Error())
	}
	time.Sleep(d)
}

func main() {
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
	}

	for i := 1; i < len(os.Args); i++ {
		if strings.ContainsAny(os.Args[i], "smh") {
			sleep(os.Args[i])
		} else {
			sleep(os.Args[i] + "s")
		}
	}
}
