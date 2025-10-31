package client

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	nspdURL    = "https://nspd.gov.ru/api/geoportal/v2/search/geoportal"
	requestURL = nspdURL + "?query=%s"
)

type NSDPClient struct {
	client *http.Client
}

func NewNSDPClient() *NSDPClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &NSDPClient{
		client: client,
	}
}

func (c *NSDPClient) GetZoneDetails(ctx context.Context, zoneID string) (*model.NSPDResp, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(requestURL, zoneID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("User-Agent", "MyGoApp/1.0")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("referer", "https://nspd.gov.ru/map?thematic=PKK&zoom=20&coordinate_x=4187280.1010340527&coordinate_y=7507815.775997361&theme_id=1&is_copy_url=true&active_layers=%E8%B3%91%2C%E8%B3%90")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	details := &model.NSPDResp{}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	//if err := json.NewDecoder(resp.Body).Decode(details); err != nil {
	//	return nil, fmt.Errorf("failed to decode response: %w", err)
	//}

	return details, nil
}
