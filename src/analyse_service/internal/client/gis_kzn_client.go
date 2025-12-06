package client

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	zoneIdURL      = "https://gis.kzn.ru/api/gis/ows/2/portal/wms?service=wms&request=GetFeatureInfo&bbox=%f,%f,%f,%f&layers=portal:portal_geo_urban10_functional_zone_exist&query_layers=portal:portal_geo_urban10_functional_zone_exist&x=%d&y=%d&width=512&height=512&srs=EPSG:3857&info_format=application/json&propertyName=key&feature_count=15&styles=portal:V_GEO_URBAN_10_FUNCTIONAL_ZONE_exist"
	zoneDetailsURL = "https://gis.kzn.ru/api/card/portal/Urban10FunctionalZone/%s"
)

type GisKznClient struct {
	client *http.Client
}

func NewGisKznClient() *GisKznClient {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	return &GisKznClient{
		client: client,
	}
}

func (c *GisKznClient) GetZoneID(ctx context.Context, zone *model.LayerInfoCoords) (*model.GisKznZoneIDResp, error) {
	resultReq := fmt.Sprintf(zoneIdURL, zone.Bbox[1].X, zone.Bbox[1].Y, zone.Bbox[0].X, zone.Bbox[0].Y, zone.Coords.X, zone.Coords.Y)

	log.Printf("Requesting: %s", resultReq)
	req, err := http.NewRequest("GET", resultReq, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("User-Agent", "MyGoApp/1.0")
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gis kzn zone details: %s", resp.Status)
	}

	details := &model.GisKznZoneIDResp{}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return details, nil
}

func (c *GisKznClient) GetFuncZoneDetails(zoneID string) (*model.FuncZoneDetailsInfo, error) {
	resultReq := fmt.Sprintf(zoneDetailsURL, zoneID)

	log.Printf("Requesting: %s", resultReq)
	req, err := http.NewRequest("GET", resultReq, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("User-Agent", "MyGoApp/1.0")
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gis kzn zone details: %s", resp.Status)
	}

	details := &model.FuncZoneDetailsInfo{}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return details, nil
}
