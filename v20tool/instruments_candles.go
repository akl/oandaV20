package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gilramir/argparse/v2"
	"github.com/gilramir/oandaV20"
)

type InstrumentsCandlesOptions struct {
	RootOptions

	Save string

	Instrument        string
	Price             string
	Granularity       string
	Count             int
	From              string
	To                string
	Smooth            bool
	NoIncludeFirst    bool
	DailyAlignment    int
	AlignmentTimezone string
	WeeklyAlignment   string
}

func build_argparse_instruments_candles(pcmd *argparse.Command) {
	ap := pcmd.New(&argparse.Command{
		Name:        "candles",
		Description: "Call the instruments candles API",
		Function:    RunInstrumentsCandles,
		Values: &InstrumentsCandlesOptions{
			RootOptions: RootOptions{
				ConfigFile: kDefaultConfigFile,
			},
		},
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--save"},
		MetaVar:  "FILE",
		Help:     "Save the raw JSON to FILE",
	})

	ap.Add(&argparse.Argument{
		Name: "instrument",
		Help: "The OANDA instrument symbol",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--price", "-p"},
		Help:     "Any combination of M (mid), B (bid) and A (ask).",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--granularity", "-g"},
		Choices:  oandaV20.ValidCandlestickGranularities,
		Help:     "The granularity of a candlestick.",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--count"},
		Help:     "The number of candlesticks to return",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--from"},
		Help:     "Start time, in RFC3339 format: 2006-01-02T15:04:05Z07:00",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--to"},
		Help: `End time, in RFC3339 format: 2006-01-02T15:04:05Z07:00
			Or, the word "now" can be given.`,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--smooth"},
		Help:     "Should the candlestick be \"smoothed\"?",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--no-include-first"},
		Help:     "Don't include the first candlestick in the response",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--daily-alignment", "-da"},
		Help: `The hour of the day to use for granularites that
		have daily alignment.`,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--alignment-timezone", "-at"},
		Help: `The timezone to be used for the daily
		alignment parameter. The returned times will still be
		represented in UTC.`,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--weekly-alignment", "-wa"},
		Help: `The day of the week used for granularities
		that have weekly alignment.`,
		Choices: []string{"Monday", "Tuesday", "Wednesday",
			"Thursday", "Friday", "Saturday", "Sunday"},
	})

}

func RunInstrumentsCandles(cmd *argparse.Command, values argparse.Values) error {
	opts := values.(*InstrumentsCandlesOptions)

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	log.Printf("Got %v", opts)

	log.Printf("opts.ConfigFile is %s", opts.ConfigFile)
	cfg, err := ReadConfig(opts.ConfigFile)
	if err != nil {
		return err
	}

	client, err := oandaV20.NewClient(cfg.AccessToken,
		cfg.Environment, cfg.Streaming)
	if err != nil {
		return err
	}

	req := oandaV20.DefaultCandleRequest()

	// Create the request
	if cmd.Seen["Price"] {
		req.Price = opts.Price
	}
	if cmd.Seen["Granularity"] {
		req.Granularity = opts.Granularity
	}
	if cmd.Seen["Count"] {
		if cmd.Seen["From"] && cmd.Seen["To"] {
			return errors.New("--count cannot be given if both " +
				"--to and --from are given.")
		}
		req.Count = opts.Count
	}
	if cmd.Seen["From"] {
		req.From, err = time.Parse(time.RFC3339, opts.From)
		log.Printf("Set From to %v", req.From)
		if err != nil {
			return err
		}
	}
	if cmd.Seen["To"] {
		if opts.To == "now" {
			req.To = time.Now()
		} else {
			req.To, err = time.Parse(time.RFC3339, opts.To)
			if err != nil {
				return err
			}
		}
	}
	if cmd.Seen["Smooth"] {
		req.Smooth = true
	}
	if cmd.Seen["NoIncludeFirst"] {
		req.IncludeFirst = false
	}
	if cmd.Seen["DailyAlignment"] {
		req.DailyAlignment = opts.DailyAlignment
	}
	if cmd.Seen["AlignmentTimezone"] {
		req.AlignmentTimezone = opts.AlignmentTimezone
	}
	if cmd.Seen["WeeklyAlignment"] {
		req.WeeklyAlignment = opts.WeeklyAlignment
	}

	// Save to a file?
	if opts.Save != "" {
		body, err := client.GetBidAskCandlesBytes(opts.Instrument, req)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(opts.Save, body, 0644)
		if err != nil {
			return err
		}
		return nil
	}

	instrumentCandles, err := client.GetInstrumentCandles(opts.Instrument, req)
	if err != nil {
		return err
	}
	fmt.Printf("Got: %+v\n", instrumentCandles)

	return nil
}
