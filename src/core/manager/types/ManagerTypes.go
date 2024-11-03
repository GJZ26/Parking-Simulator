package types

import (
	"Parking_Simulator/src/core/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

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

type Side byte

const (
	LeftSide  Side = 0
	RightSide Side = 1
)

type SlotInfo struct {
	Id       uint32
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Type     string
	Occupied bool
	Side     Side
}

type PointMap struct {
	Road             []PointData
	LeftParkingSlot  []PointData
	RightParkingSlot []PointData
	ParkingRoad      []PointData
}

type RenderData struct {
	Counter int
	Cars    []*entity.Car
}
