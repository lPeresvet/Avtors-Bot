package usecase

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"context"
	"fmt"
)

const (
	x = 0
	y = 1

	min = 0
	max = 1

	bboxSize = 512 //TODO maybe increase
)

type NSPDClient interface {
	GetTerrZone(ctx context.Context, zone *model.LayerInfoCoords) (*model.TerrZone, error)
}

type ZonesUseCase struct {
	nspdClient NSPDClient
}

func NewZonesUseCase(nspdClient NSPDClient) *ZonesUseCase {
	return &ZonesUseCase{
		nspdClient: nspdClient,
	}
}

func (uc *ZonesUseCase) GetZones(ctx context.Context, coords [][]float64) (*model.ZonesAnalysis, error) {
	if len(coords) < 3 {
		return nil, fmt.Errorf("not enough coordinates")
	}
	outCoords := outerBorder(coords[0 : len(coords)-1])

	zone, err := uc.nspdClient.GetTerrZone(ctx, &model.LayerInfoCoords{
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
	})
	if err != nil {
		return nil, err
	}

	return &model.ZonesAnalysis{
		TerrZoneName: zone.Features[0].Properties.Options.NameByDoc,
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
