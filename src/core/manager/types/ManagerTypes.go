package types

import "github.com/hajimehoshi/ebiten/v2"

type ResourcesRequest struct {
	iconPath string
}

type CarSprites struct {
	Blue   *ebiten.Image
	Green  *ebiten.Image
	Orange *ebiten.Image
	Pink   *ebiten.Image
	Red    *ebiten.Image
	Yellow *ebiten.Image
}

type CarSpritePath struct {
	Blue   string
	Green  string
	Orange string
	Pink   string
	Red    string
	Yellow string
}

type PointData struct {
	Id     uint32
	X      float64
	Y      float64
	Width  float64
	Height float64
	Type   string
}

type PointMap struct {
	Road        []PointData
	ParkingSlot []PointData
	ParkingRoad []PointData
}
