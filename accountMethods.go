package imoose

import (
	"encoding/json"
	"net/url"
)

type getPortfoliosResponse struct {
	APIResponse
	Data []Portfolio `json:"data"`
}

type getPortfolioBalancesResponse struct {
	APIResponse
	Data map[string]json.Number `json:"data"`
}

func (c Client) GetPortfolios() ([]Portfolio, error) {

	resp := getPortfoliosResponse{}

	err := c.get("/v1/private/portfolio", url.Values{}, &resp)
	if err != nil {
		return []Portfolio{}, err
	}

	return resp.Data, nil
}

func (c Client) GetPortfolioBalance(portfolioId string) (map[string]float64, error) {
	resp := getPortfolioBalancesResponse{}
	balances := map[string]float64{}
	err := c.get("/v1/private/balance", url.Values{"id": []string{portfolioId}}, &resp)
	if err != nil {
		return map[string]float64{}, err
	}
	for key, value := range resp.Data {
		balances[key], _ = value.Float64()
	}
	return balances, nil
}
