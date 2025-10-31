package adapter

import (
	"avtor.ru/bot/client"
	"context"
	"encoding/json"
	"fmt"
)

type NspdAdapter struct {
	client *client.Client
}

func NewNspd(address string) (*NspdAdapter, error) {
	nspdClient, err := client.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create nspd client: %s", err)
	}

	return &NspdAdapter{client: nspdClient}, nil
}

func (n *NspdAdapter) Analyse(ctx context.Context, zoneID string) (*client.ZoneDetails, error) {
	resp, err := n.client.GetZonesZoneIDAnalise(ctx, zoneID)
	if err != nil {
		return &client.ZoneDetails{}, fmt.Errorf("failed to get zone analise: %s", err)
	}

	zoneInfo := &client.ZoneDetails{}
	if err := json.NewDecoder(resp.Body).Decode(zoneInfo); err != nil {
		return &client.ZoneDetails{}, fmt.Errorf("failed to decode zone info: %s", err)
	}

	return zoneInfo, nil
}
