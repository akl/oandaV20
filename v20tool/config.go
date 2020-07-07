package main

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"

	"github.com/gilramir/oandaV20"
	"gopkg.in/ini.v1"
)

type Config struct {
	AccessToken string
	Environment string
	Streaming   bool
}

func ReadConfig(filename string) (*Config, error) {

	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("Reading INI file %s: %q",
			filename, err)
	}

	c := &Config{}
	c.AccessToken = cfg.Section("oanda").Key("access_token").String()
	c.Environment = cfg.Section("oanda").Key("env").String()
	c.Streaming, err = cfg.Section("oanda").Key("streaming").Bool()
	if err != nil {
		return nil, fmt.Errorf("Parsing [oanda]/streaming: %q", err)
	}

	if c.Environment != oandaV20.TradeEnvironment &&
		c.Environment != oandaV20.PracticeEnvironment {
		return nil, fmt.Errorf("[oanda]/env should be either %s or %s",
			oandaV20.TradeEnvironment, oandaV20.PracticeEnvironment)
	}

	return c, nil
}
