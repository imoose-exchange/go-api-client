package imoose

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type GetMarketsParamters struct {
	MarketType MarketType
}

type getMarketResponse struct {
	APIResponse
	Data Market `json:"data"`
}

type getMarketsResponse struct {
	APIResponse
	Data []Market `json:"data"`
}

type getTickerResponse struct {
	APIResponse
	Data Ticker `json:"data"`
}

type getTickersResponse struct {
	APIResponse
	Data []Ticker `json:"data"`
}

type getAssetResponse struct {
	APIResponse
	Data Asset `json:"data"`
}

type getAssetsResponse struct {
	APIResponse
	Data []Asset `json:"data"`
}

type getRecentMarketTradesResponse struct {
	APIResponse
	Data [][]interface{} `json:"data"`
}

type getMarketOHLCResponse struct {
	APIResponse
	Data [][]json.Number `json:"data"`
}

type getMarketDepthResponse struct {
	APIResponse
	Data struct {
		Bids [][]json.Number `json:"bids"`
		Asks [][]json.Number `json:"asks"`
	} `json:"data"`
}

func (paramters GetMarketsParamters) getValues() url.Values {
	values := url.Values{}
	if len(paramters.MarketType) <= 0 {
		values.Add("type", string(MARKET_TYPE_SPOT))
	} else {
		values.Add("type", string(paramters.MarketType))
	}
	return values
}

func (c Client) GetLiveMarkets(parameters GetMarketsParamters) ([]Market, error) {

	resp := getMarketsResponse{}

	err := c.get("/v1/public/market", parameters.getValues(), &resp)
	if err != nil {
		return []Market{}, err
	}

	return resp.Data, nil
}

func (c Client) GetMarket(id string) (Market, error) {

	resp := getMarketResponse{}

	err := c.get("/v1/public/market", url.Values{"id": []string{id}}, &resp)
	if err != nil {
		return Market{}, err
	}
	return resp.Data, nil
}

type GetMarketTickersParamters struct {
	MarketType MarketType
}

func (parameters GetMarketTickersParamters) getValues() url.Values {
	values := url.Values{}
	if len(parameters.MarketType) <= 0 {
		values.Add("type", string(MARKET_TYPE_SPOT))
	} else {
		values.Add("type", string(parameters.MarketType))
	}
	return values
}

func (c Client) GetMarketTickers(parameters GetMarketTickersParamters) ([]Ticker, error) {
	resp := getTickersResponse{}
	err := c.get("/v1/public/ticker", parameters.getValues(), &resp)
	if err != nil {
		return []Ticker{}, err
	}
	return resp.Data, nil
}

func (c Client) GetMarketTicker(id string) (Ticker, error) {
	resp := getTickerResponse{}
	err := c.get("/v1/public/ticker", url.Values{"id": []string{id}}, &resp)
	if err != nil {
		return Ticker{}, err
	}
	return resp.Data, nil
}

type GetRecentMarketTradesParamaters struct {
	ID    string
	Limit int
}

func (paramters GetRecentMarketTradesParamaters) getValues() url.Values {
	values := url.Values{}
	values.Add("id", paramters.ID)

	if paramters.Limit > 0 {
		values.Add("limit", fmt.Sprint(paramters.Limit))
	}
	return values
}

// returns array of trades [price, volume, time]
func (c Client) GetRecentMarketTrades(paramters GetRecentMarketTradesParamaters) ([]MarketTrade, error) {

	resp := &getRecentMarketTradesResponse{}

	err := c.get("/v1/public/trade", paramters.getValues(), &resp)
	if err != nil {
		return []MarketTrade{}, err
	}

	data := []MarketTrade{}

	for _, input := range resp.Data {
		trade := MarketTrade{}
		trade.Price, _ = strconv.ParseFloat(input[0].(string), 64)
		trade.Volume, _ = strconv.ParseFloat(input[1].(string), 64)
		trade.Time, _ = input[2].(int64)
		data = append(data, trade)
	}

	return data, nil
}

type GetMarketDepthParamaters struct {
	ID    string
	Limit int
}

func (paramters GetMarketDepthParamaters) getValues() url.Values {
	values := url.Values{}
	values.Add("id", paramters.ID)

	if paramters.Limit > 0 {
		values.Add("limit", fmt.Sprint(paramters.Limit))
	}
	return values
}

// returns array of trades [price, volume, time]
func (c Client) GetMarketDepth(paramters GetMarketDepthParamaters) (OrderBook, error) {

	resp := &getMarketDepthResponse{}

	err := c.get("/v1/public/depth", paramters.getValues(), &resp)
	if err != nil {
		return OrderBook{}, err
	}

	book := OrderBook{
		Asks: [][]float64{},
		Bids: [][]float64{},
	}

	for _, v := range resp.Data.Asks {
		price, _ := v[0].Float64()
		volume, _ := v[1].Float64()
		book.Asks = append(book.Asks, []float64{price, volume})
	}

	for _, v := range resp.Data.Bids {
		price, _ := v[0].Float64()
		volume, _ := v[1].Float64()
		book.Bids = append(book.Bids, []float64{price, volume})
	}

	return book, nil
}

type GetMarketOHLCParamters struct {
	ID       string
	Interval int
	Since    int64
}

func (paramters GetMarketOHLCParamters) getValues() url.Values {
	values := url.Values{}
	values.Add("id", paramters.ID)
	if paramters.Interval > 0 {
		values.Add("interval", fmt.Sprint(paramters.Interval))
	}
	if paramters.Since > 0 {
		values.Add("since", fmt.Sprint(paramters.Since))
	}
	return values
}

// returns array of trades [price, volume, time]
func (c Client) GetMarketOHLC(paramters GetMarketOHLCParamters) ([]OHLCItem, error) {

	resp := &getMarketOHLCResponse{}

	err := c.get("/v1/public/ohlc", paramters.getValues(), &resp)
	if err != nil {
		return []OHLCItem{}, err
	}

	data := []OHLCItem{}

	for _, input := range resp.Data {
		trade := OHLCItem{}
		trade.Time, _ = input[0].Int64()
		trade.Open, _ = input[1].Float64()
		trade.High, _ = input[2].Float64()
		trade.Low, _ = input[3].Float64()
		trade.Close, _ = input[4].Float64()
		trade.Volume, _ = input[5].Float64()
		data = append(data, trade)
	}

	return data, nil
}

func (c Client) GetAsset(id string) (Asset, error) {

	resp := getAssetResponse{}

	err := c.get("/v1/public/asset", url.Values{"id": []string{id}}, &resp)
	if err != nil {
		return Asset{}, err
	}

	return resp.Data, nil
}

func (c Client) GetAssets() ([]Asset, error) {

	resp := getAssetsResponse{}

	err := c.get("/v1/public/asset", url.Values{}, &resp)
	if err != nil {
		return []Asset{}, err
	}

	return resp.Data, nil
}
