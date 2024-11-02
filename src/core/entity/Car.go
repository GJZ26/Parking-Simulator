package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Car struct {
	sprite *ebiten.Image
	x      int
	y      int
	angle  float64
}

func NewCar(x int, y int, sprite *ebiten.Image) *Car {
	return &Car{x: x, y: y, angle: 0, sprite: sprite}
}
func (c *Car) Update() {
	c.x++
	c.angle++
}

func (c *Car) GetX() int {
	return c.x
}

func (c *Car) GetY() int {
	return c.y
}

func (c *Car) GetAngle() float64 {
	return c.angle
}

func (c *Car) GetSprite() *ebiten.Image {
	return c.sprite
}
