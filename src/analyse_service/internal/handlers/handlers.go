package handlers

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"avtor.ru/bot/server"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	_ "net/http"
	"time"
)

type NSPDClient interface {
	GetZoneDetails(ctx context.Context, zoneID string) (*model.NSPDResp, error)
}

type SubZonesService interface {
	GetZones(ctx context.Context, coords [][]float64) (*model.ZonesAnalysis, error)
}

type Repository interface {
	InsertLike(like model.Like) error
	DeleteLike(like model.Like) error
	GetLikes() (*server.Zones, error)
	GetUserRole(username string) (*model.Role, error)
	CreateUser(username, role string) error
	GetUsers() (*server.Users, error)
	DeleteUser(username string) error
}

type AnalyseService struct {
	ctx        context.Context
	nspdClient NSPDClient
	repo       Repository
	subZones   SubZonesService
}

func NewAnalyseService(ctx context.Context, nspdClient NSPDClient, repository Repository, subZones SubZonesService) *AnalyseService {
	return &AnalyseService{
		ctx:        ctx,
		nspdClient: nspdClient,
		repo:       repository,
		subZones:   subZones,
	}
}

func (svc *AnalyseService) GetUserZones(ctx echo.Context) error {
	zones, err := svc.repo.GetLikes()
	if err != nil {
		log.Printf("GetLikes: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{Code: http.StatusInternalServerError, Message: "Failed to get likes"})
	}

	return ctx.JSON(http.StatusOK, zones)
}

func (svc *AnalyseService) DeleteZonesZoneIDLikeUserID(ctx echo.Context, zoneID string, userID string) error {
	err := svc.repo.DeleteLike(model.Like{
		ZoneID: zoneID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Failed to unlike zone: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete zone like",
		})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (svc *AnalyseService) GetZonesZoneIDAnalise(ctx echo.Context, zoneID string) error {
	//TODO: add zone id validation and normal logs
	timeoutCtx, cancel := context.WithTimeout(svc.ctx, 20*time.Second)
	defer cancel()

	details, err := svc.nspdClient.GetZoneDetails(timeoutCtx, zoneID)
	if err != nil {
		log.Printf("GetZoneDetails: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get zone details",
		})
	}

	zones, err := svc.subZones.GetZones(timeoutCtx, details.Data.Features[0].Geometry.Coordinates[0])
	if err != nil {
		log.Printf("GetSubZoneDetails: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get sub zone details",
		})
	}

	serviceResponse := &server.ZoneDetails{
		Id:              zoneID,
		PermittedUsage:  details.Data.Features[0].Properties.Options.PermittedUseEstablishedByDocument,
		PropertyType:    ConvertOwnershipType(details.Data.Features[0].Properties.Options.OwnershipType),
		RightType:       &details.Data.Features[0].Properties.Options.RightType,
		Square:          &details.Data.Features[0].Properties.Options.LandRecordAreaVerified,
		Address:         &details.Data.Features[0].Properties.Options.ReadableAddress,
		FunctionalZone:  &zones.FunctionalZoneName,
		TerritorialZone: &zones.TerrZoneName,
		ZoneOK:          &zones.ConstructionPermitted,
	}

	return ctx.JSON(http.StatusOK, serviceResponse)
}

func (svc *AnalyseService) PostZonesZoneIDLikeUserID(ctx echo.Context, zoneID string, userID string) error {
	err := svc.repo.InsertLike(model.Like{
		ZoneID: zoneID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("InsertLike: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{Code: http.StatusInternalServerError, Message: "Failed to insert like"})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (svc *AnalyseService) GetUserAuthUserID(ctx echo.Context, username string) error {
	role, err := svc.repo.GetUserRole(username)
	if err != nil {
		log.Printf("GetUserRole: %v", err)

		return ctx.JSON(http.StatusNotFound, server.Error{Code: http.StatusNotFound, Message: "Failed to get user role"})
	}

	convertedRole := ConvertRole(*role)
	if convertedRole == server.UndefinedRole {
		log.Printf("GetUserRole failed to convert role: %v", convertedRole)

		return ctx.JSON(http.StatusNotFound, server.Error{Code: http.StatusNotFound, Message: "Failed to get user role"})
	}

	return ctx.JSON(http.StatusOK, convertedRole)
}

func (svc *AnalyseService) PostUserCreateUserID(ctx echo.Context, _ string) error {
	var user model.User
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		log.Printf("PostUserCreateUserID: %v", err)
		return ctx.JSON(http.StatusBadRequest, server.Error{Code: http.StatusBadRequest, Message: "Failed to parse request body"})
	}

	if err := svc.repo.CreateUser(user.Username, string(user.Role)); err != nil {
		log.Printf("CreateUser: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{Code: http.StatusInternalServerError, Message: "Failed to create user"})
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (svc *AnalyseService) GetUsers(ctx echo.Context) error {
	users, err := svc.repo.GetUsers()
	if err != nil {
		log.Printf("GetUsers error: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{Code: http.StatusInternalServerError, Message: "Failed to get users"})
	}

	return ctx.JSON(http.StatusOK, users)
}

func (svc *AnalyseService) DeleteUsersUserID(ctx echo.Context, userID string) error {
	if err := svc.repo.DeleteUser(userID); err != nil {
		log.Printf("GetUsers error: %v", err)

		return ctx.JSON(http.StatusInternalServerError, server.Error{Code: http.StatusInternalServerError, Message: "Failed to delete users"})
	}

	return ctx.JSON(http.StatusAccepted, nil)
}
