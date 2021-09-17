package imoose

import (
	"fmt"
	"net/url"
)

type PlaceOrderParamters struct {
	PortfolioID string
	OrderType   OrderType
	MarketID    string
	Side        OrderSide
	Price       float64
	Volume      float64
}

type placeOrderResponse struct {
	APIResponse
	Data Order `json:"data"`
}

type queryOrderResponse struct {
	APIResponse
	Data Order `json:"data"`
}

type cancelOrderResponse struct {
	APIResponse
	Data struct {
		Canceled int `json:"canceled"`
	} `json:"data"`
}

type queryOrdersResponse struct {
	APIResponse
	Data struct {
		Orders []Order `json:"orders"`
		Next   string  `json:"next"`
	} `json:"data"`
}

func (paramters PlaceOrderParamters) getValues() url.Values {
	values := url.Values{}
	values.Add("portfolio_id", paramters.PortfolioID)
	values.Add("type", string(paramters.OrderType))
	values.Add("side", string(paramters.Side))
	values.Add("volume", fmt.Sprint(paramters.Volume))
	values.Add("market", paramters.MarketID)
	if paramters.OrderType == ORDER_TYPE_LIMIT {
		values.Add("price", fmt.Sprint(paramters.Price))
	}
	return values
}

func (c Client) PlaceOrder(parameters PlaceOrderParamters) (Order, error) {
	resp := placeOrderResponse{}
	err := c.post("/v1/private/order", parameters.getValues(), &resp)
	if err != nil {
		return Order{}, err
	}
	return resp.Data, nil
}

func (c Client) QueryOrder(id string) (Order, error) {
	resp := queryOrderResponse{}
	err := c.get("/v1/private/order", url.Values{"id": []string{id}}, &resp)
	if err != nil {
		return Order{}, err
	}
	return resp.Data, nil
}

func (c Client) CancelOrder(id string) (int, error) {
	resp := cancelOrderResponse{}
	err := c.delete("/v1/private/order", url.Values{"id": []string{id}}, &resp)
	if err != nil {
		return 0, err
	}
	return resp.Data.Canceled, nil
}

type QueryOrderParamters struct {
	PortfolioID string
	Limit       int
	From        string
}

func (paramters QueryOrderParamters) getValues() url.Values {
	values := url.Values{}
	values.Add("portfolio_id", paramters.PortfolioID)
	values.Add("from", paramters.From)
	values.Add("limit", fmt.Sprint(paramters.Limit))
	return values
}

func (c Client) QueryOpenOrders(parameters QueryOrderParamters) ([]Order, string, error) {
	resp := queryOrdersResponse{}
	err := c.get("/v1/private/order/open", parameters.getValues(), &resp)
	if err != nil {
		return []Order{}, "", err
	}
	return resp.Data.Orders, resp.Data.Next, nil
}

func (c Client) QueryOrderHistory(parameters QueryOrderParamters) ([]Order, string, error) {
	resp := queryOrdersResponse{}
	err := c.get("/v1/private/order/closed", parameters.getValues(), &resp)
	if err != nil {
		return []Order{}, "", err
	}
	return resp.Data.Orders, resp.Data.Next, nil
}
