package adapter

import (
	"avtor.ru/bot/client"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var ErrorLikeZone = errors.New("failed to like zone")

type ServiceAdapter struct {
	client *client.Client
}

func NewAnalyseServiceAdapter(address string) (*ServiceAdapter, error) {
	nspdClient, err := client.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create nspd client: %s", err)
	}

	return &ServiceAdapter{client: nspdClient}, nil
}

func (n *ServiceAdapter) Analyse(ctx context.Context, zoneID string) (*client.ZoneDetails, error) {
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

func (n *ServiceAdapter) GetLikes(ctx context.Context, userID int64) (*client.Zones, error) {
	zones := &client.Zones{}
	resp, err := n.client.GetUserUserIDZones(ctx, strconv.FormatInt(userID, 10))
	if err != nil {
		return zones, fmt.Errorf("failed to get user likes: %s", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(zones); err != nil {
		return zones, fmt.Errorf("failed to decode user likes: %s", err)
	}

	return zones, nil
}

func (n *ServiceAdapter) LikeZone(ctx context.Context, userID int64, zoneID string) error {
	resp, err := n.client.PostZonesZoneIDLikeUserID(ctx, zoneID, strconv.FormatInt(userID, 10))
	if err != nil {
		return fmt.Errorf("failed to like zone: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorLikeZone
	}

	return nil
}
