package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type Car struct {
	targetX     int
	targetY     int
	sprite      *ebiten.Image
	x           int
	y           int
	angle       float64
	parkingTime int
	startTime   time.Time
	slotId      uint32
}

func NewCar(x int, y int, slotId uint32, sprite *ebiten.Image) *Car {
	return &Car{x: x, y: y, angle: 0, sprite: sprite, parkingTime: 3 + rand.Intn(3), slotId: slotId}
}

func (c *Car) Update() bool {

	if c.startTime.IsZero() {
		c.StartTimer()
	}

	return c.IsActive()
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

func (c *Car) StartTimer() {
	c.startTime = time.Now()
}

func (c *Car) IsActive() bool {
	if c.startTime.IsZero() {
		return true
	}
	endTime := c.startTime.Add(time.Duration(c.parkingTime) * time.Second)
	return time.Now().Before(endTime)
}

func (c *Car) GetSlotID() uint32 {
	return c.slotId
}
