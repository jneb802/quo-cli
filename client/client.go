package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://api.openphone.com/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type Message struct {
	ID            string   `json:"id"`
	To            []string `json:"to"`
	From          string   `json:"from"`
	Text          string   `json:"text"`
	PhoneNumberID string   `json:"phoneNumberId"`
	Direction     string   `json:"direction"`
	UserID        string   `json:"userId"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
}

type SendRequest struct {
	Content string   `json:"content"`
	From    string   `json:"from"`
	To      []string `json:"to"`
}

type ListResponse struct {
	Data          []Message `json:"data"`
	TotalItems    int       `json:"totalItems"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

type MessageResponse struct {
	Data Message `json:"data"`
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Body)
}

func (c *Client) do(req *http.Request, out interface{}) error {
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("parsing response: %w", err)
		}
	}
	return nil
}

func (c *Client) SendMessage(from, to, content string) (*Message, error) {
	payload, err := json.Marshal(SendRequest{
		Content: content,
		From:    from,
		To:      []string{to},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+"/messages", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var result MessageResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) ListMessages(phoneNumberID, participant string, maxResults int, pageToken string) (*ListResponse, error) {
	params := url.Values{}
	params.Set("phoneNumberId", phoneNumberID)
	params.Add("participants[]", participant)
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}

	req, err := http.NewRequest("GET", baseURL+"/messages?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetMessage(id string) (*Message, error) {
	req, err := http.NewRequest("GET", baseURL+"/messages/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}

	var result MessageResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
