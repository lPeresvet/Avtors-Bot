package handlers

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"avtor.ru/bot/server"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	_ "net/http"
	"time"
)

type NSPDClient interface {
	GetZoneDetails(ctx context.Context, zoneID string) (*model.NSPDResp, error)
}

type AnalyseService struct {
	ctx        context.Context
	nspdClient NSPDClient
}

func NewAnalyseService(ctx context.Context, nspdClient NSPDClient) *AnalyseService {
	return &AnalyseService{
		ctx:        ctx,
		nspdClient: nspdClient,
	}
}

func (svc *AnalyseService) GetUserUserIDZones(ctx echo.Context, userID string) error {
	return nil
}

func (svc *AnalyseService) GetZonesZoneIDAnalise(ctx echo.Context, zoneID string) error {
	//TODO: add zone id validation and normal logs
	timeoutCtx, cancel := context.WithTimeout(svc.ctx, 5*time.Second)
	defer cancel()

	details, err := svc.nspdClient.GetZoneDetails(timeoutCtx, zoneID)
	if err != nil {
		log.Printf("GetZoneDetails: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get zone details",
		})
	}

	serviceResponse := &server.ZoneDetails{
		Id:             zoneID,
		PermittedUsage: details.Data.Features[0].Properties.Options.PermittedUseEstablishedByDocument,
		PropertyType:   ConvertOwnershipType(details.Data.Features[0].Properties.Options.OwnershipType),
		RightType:      &details.Data.Features[0].Properties.Options.RightType,
	}

	return ctx.JSON(http.StatusOK, serviceResponse)
}

func (svc *AnalyseService) PostZonesZoneIDLike(ctx echo.Context, zoneID string) error {
	return nil
}
