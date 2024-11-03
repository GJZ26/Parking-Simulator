package entity

import (
	"Parking_Simulator/src/core/manager/types/geo"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

const carSpeed = 3

type Car struct {
	targetSlot    geo.SlotInfo
	sprite        *ebiten.Image
	x             int
	y             int
	angle         float64
	parkingTime   int
	startTime     time.Time
	defaultRoad   []geo.PointData
	parkingRoad   []geo.PointData
	queuePosition int
}

func NewCar(target geo.SlotInfo, road []geo.PointData, parking []geo.PointData, sprite *ebiten.Image, queuePosition int) *Car {
	return &Car{
		x:             int(road[0].X),
		y:             int(road[0].Y),
		angle:         -90,
		sprite:        sprite,
		parkingTime:   3 + rand.Intn(3),
		targetSlot:    target,
		parkingRoad:   parking,
		defaultRoad:   road,
		queuePosition: queuePosition,
	}
}

func (c *Car) Update() bool {

	if c.x >
		int(c.defaultRoad[1].X+ // posición x de la ruta principal
			90+ // Distancia desde la intersacción de la entrada
			(70*float64(c.queuePosition))) { // Distancia entre auto en cola
		c.x -= carSpeed
	}

	c.StartTimer()
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
	return c.targetSlot.Id
}

func (c *Car) SetQueuePosition(position int) {
	c.queuePosition = position
}
