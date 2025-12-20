package usecase

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	x = 0
	y = 1

	min = 0
	max = 1

	bboxSize = 512 //TODO maybe increase
)

var (
	failedToGetZone = "Зона неизвестна"
	falseVal        = false
)

var (
	permittedFuncZones = map[string]bool{
		"120": true,
		"123": true,
		"124": true,
		"200": true,
		"202": true,
		"203": true,
		"204": true,
		"210": true,
		"230": true,
		"320": true,
	}

	permittedTerrZones = map[string]bool{
		"Зона индивидуальной жилой застройки в границах многоквартирной застройки": true,
		"Зона малоэтажной жилой застройки":                                         true,
		"Зона среднеэтажной жилой застройки":                                       true,
		"Зона многоэтажной жилой застройки":                                        true,
		"Зона смешанной жилой и общественной застройки":                            true,
		"Зона смешанной застройки центра города":                                   true,
		"Зона фактического использования территории":                               true,
	}
)

type NSPDClient interface {
	GetTerrZone(ctx context.Context, zone *model.LayerInfoCoords) (*model.TerrZone, error)
}

type GisKznClient interface {
	GetZoneID(ctx context.Context, zone *model.LayerInfoCoords) (*model.GisKznZoneIDResp, error)
	GetFuncZoneDetails(zoneID string) (*model.FuncZoneDetailsInfo, error)
}

type ZonesUseCase struct {
	nspdClient   NSPDClient
	gisKznClient GisKznClient
}

func NewZonesUseCase(nspdClient NSPDClient, gisKznClient GisKznClient) *ZonesUseCase {
	return &ZonesUseCase{
		nspdClient:   nspdClient,
		gisKznClient: gisKznClient,
	}
}

func (uc *ZonesUseCase) GetZones(ctx context.Context, coords [][]float64) (*model.ZonesAnalysis, error) {
	if len(coords) < 3 {
		return nil, fmt.Errorf("not enough coordinates")
	}
	outCoords := outerBorder(coords[0 : len(coords)-1])

	zones := &model.LayerInfoCoords{
		Bbox: []model.Coords{
			{
				X: coords[min][x],
				Y: coords[min][y],
			},
			{
				X: coords[max][x],
				Y: coords[max][y],
			},
		},
		Coords: getSubCoors(outCoords, coords[0]),
	}

	zone, err := uc.nspdClient.GetTerrZone(ctx, zones)
	if err != nil {
		return nil, err
	}

	funcZoneID, err := uc.gisKznClient.GetZoneID(ctx, zones)
	if err != nil {
		return nil, fmt.Errorf("failed to get functional zone id: %w", err)
	}

	if len(funcZoneID.Features) < 1 {
		return &model.ZonesAnalysis{
			TerrZoneName:          failedToGetZone,
			FunctionalZoneName:    failedToGetZone,
			ConstructionPermitted: false,
		}, nil
	}

	funcZoneName, err := uc.gisKznClient.GetFuncZoneDetails(strconv.FormatInt(funcZoneID.Features[0].Properties.Key, 10))
	if err != nil {
		return nil, fmt.Errorf("failed to get func zone name: %w", err)
	}

	code, err := funcZoneName.GetZoneCode()
	if err != nil {
		return nil, fmt.Errorf("failed to get func zone code: %w", err)
	}

	terrZoneName := zone.Features[0].Properties.Options.NameByDoc

	return &model.ZonesAnalysis{
		TerrZoneName:          terrZoneName,
		FunctionalZoneName:    funcZoneName.Title,
		ConstructionPermitted: isConstructionPermitted(code, terrZoneName),
	}, nil
}

func outerBorder(coords [][]float64) [][]float64 {
	minX, maxX := coords[0][x], coords[0][x]
	minY, maxY := coords[0][y], coords[0][y]

	for _, v := range coords {
		if v[x] < minX {
			minX = v[x]
		}
		if v[x] > maxX {
			maxX = v[x]
		}
		if v[y] < minY {
			minY = v[y]
		}
		if v[y] > maxY {
			maxY = v[y]
		}
	}

	return [][]float64{{minX, minY}, {maxX, maxY}}
}

func getSubCoors(borders [][]float64, point []float64) model.CoordsInt {
	xRange := borders[max][x] - borders[min][x]
	yRange := borders[max][y] - borders[min][y]

	xInterval := xRange / bboxSize
	yInterval := yRange / bboxSize

	pointX := point[x] - borders[min][x]
	pointY := point[y] - borders[min][y]

	//todo maybe fix convertion
	return model.CoordsInt{
		X: int(xInterval * pointX),
		Y: int(yInterval * pointY),
	}
}

func isConstructionPermitted(funcZoneCode, terrZoneCode string) bool {
	return isFuncZonePermitted(funcZoneCode) && isTerrZonePermitted(terrZoneCode)
}

func isFuncZonePermitted(funcZoneCode string) bool {
	log.Printf("isFuncZonePermitted(%s)", funcZoneCode)

	return permittedFuncZones[funcZoneCode]
}

func isTerrZonePermitted(terrZoneCode string) bool {
	log.Printf("isTerrZonePermitted(%s)", terrZoneCode)

	result := false
	for zone, _ := range permittedTerrZones {
		if strings.Contains(terrZoneCode, zone) {
			result = true

			break
		}
	}

	return result
}
