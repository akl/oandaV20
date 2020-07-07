package main

import (
	"github.com/gilramir/argparse/v2"
)

func build_argparse_instruments(ap *argparse.ArgumentParser) {
	iap := ap.New(&argparse.Command{
		Name:        "instruments",
		Description: "Call the instruments API",
		Values: &RootOptions{
			ConfigFile: kDefaultConfigFile,
		},
	})

	build_argparse_instruments_candles(iap)
}
