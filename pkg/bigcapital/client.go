package bigcapital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	OrgID      string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey, orgID string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		OrgID:   orgID,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, path), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("X-Organization-Id", c.OrgID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) CreateJournalEntry(entry *JournalEntry) error {
	_, err := c.doRequest("POST", "/api/v1/journal-entries", entry)
	return err
}

func (c *Client) SyncItem(item *Item) error {
	_, err := c.doRequest("POST", "/api/v1/items", item)
	return err
}

func (c *Client) SyncCustomer(customer *Customer) error {
	_, err := c.doRequest("POST", "/api/v1/customers", customer)
	return err
}
