package model

type Coords struct {
	X float64
	Y float64
}

type CoordsInt struct {
	X int
	Y int
}

type LayerInfoCoords struct {
	Bbox   []Coords
	Coords CoordsInt
}

type ZonesAnalysis struct {
	FunctionalZoneName string
	TerrZoneName       string
}
