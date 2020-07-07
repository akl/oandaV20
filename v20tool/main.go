// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"
	"github.com/gilramir/argparse/v2"
)

type RootOptions struct {
	Verbose    bool
	ConfigFile string
}

const kDefaultConfigFile = "config.ini"

func build_argparse() *argparse.ArgumentParser {
	opts := &RootOptions{
		ConfigFile: kDefaultConfigFile,
	}
	ap := argparse.New(&argparse.Command{
		Description: "The tool for testing oandaV20",
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
		Inherit:  true,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-c", "--config"},
		Dest:     "ConfigFile",
		Help: fmt.Sprintf(
			"Read this config INI file. Default is %s",
			opts.ConfigFile),
		Inherit: true,
	})

	build_argparse_instruments(ap)
	return ap
}

func main() {
	ap := build_argparse()

	ap.Parse()
}
