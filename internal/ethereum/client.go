package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kyle8615/ethereum-parser/v1/internal/model"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // expire time
		},
	}
}

func (c *Client) sendRequest(method string, params []interface{}) (*model.JSONRPCResponse, error) {
	req := model.JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// Check if the error is a timeout
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			// Handle the timeout case
			return nil, fmt.Errorf("request timed out: %v", err)
		}
		// Other errors
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var response model.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("json rpc error (code %d): %s", response.Error.Code, response.Error.Message)
	}

	return &response, nil
}

func (c *Client) GetLatestBlockNumber() (int, error) {
	response, err := c.sendRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		return -1, err
	}

	var blockNumber string
	err = json.Unmarshal(response.Result, &blockNumber)
	if err != nil {
		return -1, err
	}

	cleanedHexStr := strings.TrimPrefix(blockNumber, "0x")
	result, err := strconv.ParseInt(cleanedHexStr, 16, 64)
	if err != nil {
		return -1, err
	}
	return int(result), nil
}

func (c *Client) GetBlockByNumber(blockNumber int) (*model.Block, error) {
	response, err := c.sendRequest(
		"eth_getBlockByNumber",
		[]interface{}{
			fmt.Sprintf("0x%x", blockNumber),
			true,
		},
	)
	if err != nil {
		return nil, err
	}

	var block model.Block
	err = json.Unmarshal(response.Result, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}
