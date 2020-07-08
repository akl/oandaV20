package oandaV20

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const kPracticeRestUrl = "https://api-fxpractice.oanda.com/"
const kTradeRestUrl = "https://api-fxtrade.oanda.com/"

const kPracticeStreamingUrl = "https://stream-fxpractice.oanda.com/"
const kTradeStreamingUrl = "https://stream-fxtrade.oanda.com/"

const PracticeEnvironment = "fxpractice"
const TradeEnvironment = "fxtrade"

type Client struct {
	url string

	streaming bool

	Token string

	StreamChunkSize      int
	StreamTimeoutSeconds int

	Agent string

	// "RFC3339" or "UNIX"
	DatetimeFormat string

	client http.Client
}

/* streaming q&a
https://stackoverflow.com/questions/42959165/stream-pricing-from-oanda-v20-rest-api-using-python-requests
*/

func NewClient(accessToken string, environment string,
	streaming bool, timeoutSeconds int, keepAliveSeconds int) (*Client, error) {

	client := &Client{
		Token:                accessToken,
		StreamChunkSize:      512,
		StreamTimeoutSeconds: 10,
		Agent:                "oandaV20(go)/0.0.0",
		DatetimeFormat:       "RFC3339",
		client: http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   time.Duration(timeoutSeconds) * time.Second,
					KeepAlive: time.Duration(keepAliveSeconds) * time.Second,
				}).Dial,
				DisableKeepAlives: false,
			},
		},
	}
	if environment == PracticeEnvironment {
		if streaming {
			client.url = kPracticeStreamingUrl
		} else {
			client.url = kPracticeRestUrl
		}
	} else if environment == TradeEnvironment {
		if streaming {
			client.url = kTradeStreamingUrl
		} else {
			client.url = kTradeRestUrl
		}
	} else {
		return nil, fmt.Errorf("The 'environment' string should "+
			"be either '%s' or '%s'.", TradeEnvironment,
			PracticeEnvironment)
	}

	return client, nil
}

func (s *Client) GET(urlPath string) ([]byte, error) {

	url := s.url + urlPath
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// The v20 python code shows OANDA-Agent here
	req.Header.Set("User-Agent", s.Agent)
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Accept-Datetime-Format", "RFC3339")
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// Looks like the body will have data too. (400)
		return nil, fmt.Errorf("GET %s : %s", url, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
