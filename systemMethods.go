package imoose

import (
	"net/url"
)

type ServerStatus struct {
	Status string `json:"status"`
}

type ServerTime struct {
	UnixTime   int64  `json:"unix_time"`
	UnixTimeMS int64  `json:"unix_time_ms"`
	RFC1123    string `json:"rfc1123"`
}

type GetServerStatusResponse struct {
	APIResponse
	Data ServerStatus `json:"data"`
}

type GetServerTimerResponse struct {
	APIResponse
	Data ServerTime `json:"data"`
}

func (c Client) GetServerTime() (ServerTime, error) {
	resp := GetServerTimerResponse{}
	err := c.get("/v1/public/time", url.Values{}, &resp)
	if err != nil {
		return ServerTime{}, err
	}
	return resp.Data, nil
}

func (c Client) GetServerStatus() (ServerStatus, error) {
	resp := GetServerStatusResponse{}
	err := c.get("/v1/public/server", url.Values{}, &resp)
	if err != nil {
		return ServerStatus{}, err
	}
	return resp.Data, nil
}
