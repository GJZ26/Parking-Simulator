package resources

import (
	"Parking_Simulator/src/core/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type Request struct {
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

type RenderData struct {
	Counter int
	Cars    []*entity.Car
}
