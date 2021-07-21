package main

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type CurrencyDataWrapper struct {
	Valute	map[string]CurrencyData	`json:"Valute"`
}

type CurrencyData struct {
	Name		string	`json:"Name"`
	CharCode	string	`json:"CharCode"`
	Value		float64	`json:"Value"` // TODO decimal?
	Nominal		int		`json:"Nominal"`
}

type CurrencyBase struct {
	currencyData	map[string]CurrencyData
	url				string
}

// create struct, fill data and run worker for refilling
func InitCurrencyBase(url string, retries int, retryTimeout, timeout time.Duration) (CurrencyBase, error){
	b := CurrencyBase{
		url:          url,
	}
	err := b.loadData()

	go b.UpdateDataWorker(retries, retryTimeout, timeout)

	return b, err
}

// one try load and fill data
func (b *CurrencyBase) loadData() error {
	// load data from url
	resp, err := http.Get(b.url)
	if err != nil {
		logrus.Errorf("unable to get data: %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read data: %v", err)
		return err
	}

	// parse data
	result := CurrencyDataWrapper{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Errorf("unable to parse data: %v", err)
		return err
	}

	// save data to app
	b.currencyData = result.Valute

	return nil
}

// infinite worker for update data with timeouts
func (b *CurrencyBase) UpdateDataWorker(retries int, retryTimeout, timeout time.Duration) {
	for {
		time.Sleep(timeout)
		for tries := retries; tries > 0; tries-- {
			err := b.loadData()
			if err == nil {
				break
			}
			time.Sleep(retryTimeout)
			// TODO if not success - may be must be added a channel for error reporting and app stop ?
		}
	}
}

func (b *CurrencyBase) GetCurrency(currency string) (*CurrencyData, error) {
	item, ok := b.currencyData[currency]
	if !ok {
		return nil, errors.New("unknown currency")
	}
	return &item, nil
}

// not fast O(N), but safe get random elem from map
func (b *CurrencyBase) GetRandomCurrency() (*CurrencyData, error) {
	rnd := rand.Intn(len(b.currencyData))
	for key, _ := range b.currencyData {
		if rnd == 0 {
			return b.GetCurrency(key)
		}
		rnd--
	}
	return nil, errors.New("unreachable error")
}