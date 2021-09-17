package imoose

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	APIBaseURL   = "https://api.imoose.com" //https://api.imoose.com"
	APIUserAgent = "imoose GO API Client (https://github.com/github.com/imoose-exchange/go-api-client)"
)

type MarketType string
type OrderType string
type OrderSide string

const (
	MARKET_TYPE_SPOT    MarketType = "spot"
	MARKET_TYPE_VIRTUAL MarketType = "virtual"
	ORDER_SIDE_BUY      OrderSide  = "buy"
	ORDER_SIDE_SELL     OrderSide  = "sell"
	ORDER_TYPE_LIMIT    OrderType  = "limit"
	ORDER_TYPE_MARKET   OrderType  = "market"
)

func ParseOrderSide(side string) OrderSide {
	switch side {
	case string(ORDER_SIDE_BUY):
		return ORDER_SIDE_BUY
	case string(ORDER_SIDE_SELL):
		return ORDER_SIDE_SELL
	}
	return ""
}

func ParseOrderType(oType string) OrderType {
	switch oType {
	case string(ORDER_TYPE_LIMIT):
		return ORDER_TYPE_LIMIT
	case string(ORDER_TYPE_MARKET):
		return ORDER_TYPE_MARKET
	}
	return ""
}

func ParseMarketType(mType string) MarketType {
	switch mType {
	case string(MARKET_TYPE_SPOT):
		return MARKET_TYPE_SPOT
	case string(MARKET_TYPE_VIRTUAL):
		return MARKET_TYPE_VIRTUAL
	}
	return ""
}

// Client is the struct from which all API requests are made
type Client struct {
	key        string
	secret     string
	httpClient *http.Client
}

// ApiKeyClient instantiates the client with ApiKey Authentication
func ApiClient(key string, secret string) Client {
	return Client{
		key:        key,
		secret:     secret,
		httpClient: http.DefaultClient,
	}
}

func (c Client) WithHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

type APIResponse struct {
	Errors []string `json:"errors"`
	Status int      `json:"status"`
}

func (c Client) request(method string, endpoint string, params url.Values, holder interface{}) error {

	request, err := c.createRequest(method, endpoint, params)
	if err != nil {
		return err
	}

	var data []byte

	data, respErr := c.executeRequest(request)

	if respErr != nil {
		if len(data) > 0 {
			fmt.Println(string(data))
			imooseError := &ImooseError{}
			if err := json.Unmarshal(data, &imooseError); err != nil {
				return respErr
			}
			fmt.Println(imooseError.Errors)
			return imooseError
		}
		return respErr
	}

	if err := json.Unmarshal(data, &holder); err != nil {
		return err
	}

	return nil
}

// CreateRequest formats a request with all the necessary headers
func (c Client) createRequest(method string, endpoint string, params url.Values) (req *http.Request, err error) {
	endpoint = APIBaseURL + endpoint
	// get current unix time in MS
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	// add required timestamp value
	params.Add("timestamp", fmt.Sprint(ts))
	fmt.Println(method)
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, fmt.Sprintf("%s?%s", endpoint, params.Encode()), nil)
	} else {
		req, err = http.NewRequest(method, endpoint, strings.NewReader(params.Encode()))
		req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
		fmt.Println("body posted")
	}
	if err != nil {
		return nil, err
	}
	if strings.Contains(endpoint, "/private") {
		c.authenticate(req, params)
	}
	fmt.Println(req.URL.Path)
	req.Header.Set("User-Agent", APIUserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c Client) authenticate(req *http.Request, params url.Values) {

	plainSignature := params.Encode()

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(c.secret))
	// Write plain text signature Data to it
	h.Write([]byte(plainSignature))
	// Get hashed signature and encode as hexadecimal string
	signature := hex.EncodeToString(h.Sum(nil))

	// set authentication headers
	req.Header.Set("API-Key", c.key)
	req.Header.Set("API-Sign", signature)

	fmt.Println(c.key)
	fmt.Println(signature)
}

// executeRequest takes a prepared http.Request and returns the body of the response
// If the response is not of HTTP Code 200, an error is returned
func (c Client) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bytes := buf.Bytes()
	if resp.StatusCode != 200 {
		return bytes, fmt.Errorf("%s %s failed. Response code was %s", req.Method, req.URL, resp.Status)
	}
	return bytes, nil
}

// Get sends a GET request and marshals response data into holder
func (c Client) get(path string, params url.Values, holder interface{}) error {
	return c.request("GET", path, params, &holder)
}

// Post sends a POST request and marshals response data into holder
func (c Client) post(path string, params url.Values, holder interface{}) error {
	return c.request("POST", path, params, &holder)
}

// Delete sends a DELETE request and marshals response data into holder
func (c Client) delete(path string, params url.Values, holder interface{}) error {
	return c.request("DELETE", path, params, &holder)
}

// Put sends a PUT request and marshals response data into holder
func (c Client) put(path string, params url.Values, holder interface{}) error {
	return c.request("PUT", path, params, &holder)
}
