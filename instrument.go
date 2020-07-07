package oandaV20

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/instrument-ep/

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"
)

// From
// http://developer.oanda.com/rest-live-v20/instrument-df/#CandlestickGranularity
var ValidCandlestickGranularities []string = []string{
	"S5",  //	5 second candlesticks, minute alignment
	"S10", //	10 second candlesticks, minute alignment
	"S15", //	15 second candlesticks, minute alignment
	"S30", //	30 second candlesticks, minute alignment
	"M1",  //	1 minute candlesticks, minute alignment
	"M2",  //	2 minute candlesticks, hour alignment
	"M4",  //	4 minute candlesticks, hour alignment
	"M5",  //	5 minute candlesticks, hour alignment
	"M10", //	10 minute candlesticks, hour alignment
	"M15", //	15 minute candlesticks, hour alignment
	"M30", //	30 minute candlesticks, hour alignment
	"H1",  //	1 hour candlesticks, hour alignment
	"H2",  //	2 hour candlesticks, day alignment
	"H3",  //	3 hour candlesticks, day alignment
	"H4",  //	4 hour candlesticks, day alignment
	"H6",  //	6 hour candlesticks, day alignment
	"H8",  //	8 hour candlesticks, day alignment
	"H12", //	12 hour candlesticks, day alignment
	"D",   //	1 day candlesticks, day alignment
	"W",   //	1 week candlesticks, aligned to start of week
	"M",   //	1 month candlesticks, aligned to first day of the month
}

// Thee JSON hints for Candle, Candles, and InstrumentCandles
// were created by https://github.com/AwolDes/goanda
// (but InstrumentCandles was called InstrumentHistory)

type Candle struct {
	Open  float64 `json:"o,string"`
	Close float64 `json:"c,string"`
	Low   float64 `json:"l,string"`
	High  float64 `json:"h,string"`
}

type Candles struct {
	Complete bool      `json:"complete"`
	Volume   int       `json:"volume"`
	Time     time.Time `json:"time"`
	Mid      Candle    `json:"mid"`
	Bid      Candle    `json:"bid"`
	Ask      Candle    `json:"ask"`
}

type InstrumentCandles struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}

type CandleRequest struct {
	Price             string
	Granularity       string
	Count             int
	From              time.Time
	To                time.Time
	Smooth            bool
	IncludeFirst      bool
	DailyAlignment    int
	AlignmentTimezone string
	WeeklyAlignment   string
}

func DefaultCandleRequest() *CandleRequest {
	return &CandleRequest{
		Price:             "M",
		Granularity:       "S5",
		Count:             500,
		Smooth:            false,
		IncludeFirst:      true,
		DailyAlignment:    17,
		AlignmentTimezone: "America/New_York",
		WeeklyAlignment:   "Friday",
	}
}

func (s *Client) GetInstrumentCandles(instrument string, cr *CandleRequest) (*InstrumentCandles, error) {

	body, err := s.GetInstrumentCandlesBytes(instrument, cr)

	if err != nil {
		return nil, err
	}
	log.Printf("response %s", string(body))

	response := &InstrumentCandles{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Client) GetInstrumentCandlesBytes(instrument string, cr *CandleRequest) ([]byte, error) {
	if cr == nil {
		cr = DefaultCandleRequest()
	}
	path := s.makeCandlesUrl(instrument, cr)

	body, err := s.GET(path)
	return body, err
}

func (s *Client) makeCandlesUrl(instrument string, cr *CandleRequest) string {
	pv := url.Values{}

	pv.Set("price", cr.Price)
	pv.Set("granularity", cr.Granularity)
	if !cr.From.IsZero() {
		pv.Set("from", cr.From.Format(time.RFC3339Nano))
	}
	if !cr.To.IsZero() {
		pv.Set("to", cr.To.Format(time.RFC3339Nano))
	}
	if cr.From.IsZero() || cr.To.IsZero() {
		pv.Set("count", strconv.FormatInt(int64(cr.Count), 10))
	}

	pv.Set("smooth", strconv.FormatBool(cr.Smooth))
	pv.Set("includeFirst", strconv.FormatBool(cr.IncludeFirst))
	pv.Set("dailyAlignment", strconv.FormatInt(int64(cr.DailyAlignment), 10))
	pv.Set("alignmentTimezone", cr.AlignmentTimezone)
	pv.Set("weeklyAlignment", cr.WeeklyAlignment)

	path := "v3/instruments/" + instrument + "/candles?"

	path += pv.Encode()
	log.Printf("Path: %s", path)
	return path
}
