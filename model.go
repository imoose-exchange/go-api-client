package imoose

type Asset struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Precision int    `json:"precision"`
	State     string `json:"state"`
	Class     string `json:"class"`
	SubClass  string `json:"subclass"`
}

type Market struct {
	ID             string  `json:"id" `
	Name           string  `json:"name"`
	BaseID         string  `json:"base_id"`
	QuoteID        string  `json:"quote_id"`
	BaseName       string  `json:"base_name"`
	QuoteName      string  `json:"qoute_name"`
	Type           string  `json:"type"`
	BasePrecision  int     `json:"base_precision" `
	QuotePrecision int     `json:"quote_precision"`
	Status         string  `json:"status"`
	BuyMinVolume   float64 `json:"buy_min_volume,string"`
	BuyMaxVolume   float64 `json:"buy_max_volume,string"`
	SellMinVolume  float64 `json:"sell_min_volume,string"`
	SellMaxVolume  float64 `json:"sell_max_volume,string"`
}

type Ticker struct {
	MarketId string  `json:"market_id"`
	Open     float64 `json:"open,string"`
	Volume   float64 `json:"volume,string"`
	Low      float64 `json:"low,string"`
	High     float64 `json:"high,string"`
	Last     float64 `json:"last,string"`
	Sell     float64 `json:"sell,string"`
	Buy      float64 `json:"buy,string"`
}

type OrderBook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
}

type MarketTrade struct {
	Price  float64 `json:"p,string"`
	Volume float64 `json:"v,string"`
	Time   int64   `json:"t"`
}

type OHLCItem struct {
	Time   int64   `json:"t"`
	Open   float64 `json:"o,string"`
	High   float64 `json:"h,string"`
	Low    float64 `json:"l,string"`
	Close  float64 `json:"c,string"`
	Volume float64 `json:"v,string"`
}

type Portfolio struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Name   string `json:"name"`
}

type Order struct {
	ID               string   `json:"id"`
	Market           string   `json:"market_id"`
	Side             string   `json:"side"`
	State            string   `json:"state"`
	Type             string   `json:"type"`
	CreatedAt        int64    `json:"created_at"`
	UpdatedAt        int64    `json:"updated_at"`
	Price            float64  `json:"price,string"`
	Volume           float64  `json:"volume,string"`
	PortfolioID      string   `json:"portfolio_id"`
	Balance          float64  `json:"balance,string"`
	BalanceOrigional float64  `json:"balance_origional,string"`
	Received         float64  `json:"received,string"`
	TradeCount       int32    `json:"trade_count"`
	Fee              float64  `json:"fee,string"`
	Trades           []string `json:"trades" dynamodbav:"trades"`
}
