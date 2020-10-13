package locations

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var LocationURL = "https://www.tesla.com/all-locations"

type ClientInterface interface {
	All(types []LocationType) ([]Location, error)
	Countries(countries []string, types []LocationType) ([]Location, error)
}

type client struct {
	userAgent string
}

func NewClient(userAgent string) *client {
	return &client{
		userAgent: userAgent,
	}
}

func (c *client) All(types []LocationType) ([]Location, error) {
	var (
		all []Location
		err error
		jsonBytes []byte
	)

	cacheFile := "cache/all-locations.json"

	if stat, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		if time.Now().Sub(stat.ModTime()).Minutes() < 10 {
			jsonBytes, err = ioutil.ReadFile("cache/all-locations.json")
		}
	}

	if len(jsonBytes) == 0 {
		httpClient := http.Client{
			Timeout: 10 * time.Second,
		}

		log.Debugf("Fetching %s...", LocationURL)
		req, err := http.NewRequest(http.MethodGet, LocationURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", c.userAgent)

		res, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		jsonBytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		_ = ioutil.WriteFile(cacheFile, jsonBytes, os.ModePerm)
	}

	if err := json.Unmarshal(jsonBytes, &all); err != nil {
		return nil, err
	}

	var out []Location
	for _, l := range all {
		if l.Has(types) {
			out = append(out, l)
		}
	}

	log.Debugf("Found %d locations of type %v.", len(out), types)

	return out, err
}

func (c *client) Countries(countries []string, types []LocationType) ([]Location, error) {
	scs, err := c.All(types)
	if err != nil {
		return nil, err
	}

	var out []Location

	for _, sc := range scs {
		for _, c := range countries {
			if sc.Country == c {
				out = append(out, sc)
				break
			}
		}
	}

	return out, nil
}
