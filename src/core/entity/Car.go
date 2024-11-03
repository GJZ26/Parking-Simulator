package entity

import (
	"Parking_Simulator/src/core/manager/types/geo"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
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
	step          uint8
	internalStep  uint8
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
		step:          0,
		internalStep:  0,
	}
}

func (c *Car) Update() bool {

	switch c.internalStep {
	case 0:
		c.posInEntranceQueue()
	case 1:
		c.goToEntrance()
	case 2:
		c.goToCross()
	case 3:
		c.turnToSideSlot()
	case 4:
		c.turnUpperSideSlot()
	case 5:
		c.turnExitSideSlot()
	}

	return c.IsActive()
}

func (c *Car) posInEntranceQueue() {

	targetPosition := int(c.defaultRoad[1].X + 90 + (70 * float64(c.queuePosition)))
	isBeyondTarget := c.x > targetPosition

	if isBeyondTarget {
		c.x -= carSpeed
	} else {
		c.internalStep = 1
	}
}

func (c *Car) goToEntrance() {
	targetPosition := int(c.defaultRoad[1].X)
	isBeyondTarget := c.x > targetPosition
	isStepValid := c.step >= 1

	if isBeyondTarget && isStepValid {
		c.x -= carSpeed
	} else {
		if c.angle != 0 {
			c.angle = 0
		}
		c.internalStep = 2
	}
}

func (c *Car) goToCross() {
	targetPosition := int(c.parkingRoad[0].Y)
	isUnderTarget := c.y > targetPosition
	isStepValid := c.step >= 1

	if isUnderTarget && isStepValid {
		c.y -= carSpeed
	} else {
		c.internalStep = 3
	}
}

func (c *Car) turnToSideSlot() {
	var targetPosition float64
	var isOnTarget bool
	var angle float64
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		targetPosition = c.parkingRoad[1].X
		isOnTarget = c.x <= int(targetPosition)
		angle = -90
		speed = -(carSpeed)
	} else {
		targetPosition = c.parkingRoad[3].X
		isOnTarget = c.x >= int(targetPosition)
		angle = 90
		speed = carSpeed
	}

	if c.tryPark() {
		return
	}

	isStepValid := c.step >= 1
	c.angle = angle

	if !isOnTarget && isStepValid {
		c.x += speed
	} else {
		c.internalStep = 4
	}
}

func (c *Car) turnUpperSideSlot() {
	var targetPosition float64
	var isOnTarget bool
	var angle float64
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		targetPosition = c.parkingRoad[2].Y
		isOnTarget = c.y <= int(targetPosition)
		angle = 0
		speed = -(carSpeed)
	} else {
		targetPosition = c.parkingRoad[4].Y
		isOnTarget = c.y <= int(targetPosition)
		angle = 0
		speed = -(carSpeed)
	}

	if c.tryPark() {
		return
	}

	isStepValid := c.step >= 1
	c.angle = angle

	if !isOnTarget && isStepValid {
		c.y += speed
	} else {
		c.internalStep = 5
	}
}

func (c *Car) turnExitSideSlot() {
	var targetPosition float64
	var isOnTarget bool
	var angle float64
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		targetPosition = c.parkingRoad[5].X
		isOnTarget = c.x >= int(targetPosition)
		angle = 90
		speed = carSpeed
	} else {
		targetPosition = c.parkingRoad[5].X
		isOnTarget = c.x <= int(targetPosition)
		angle = -90
		speed = -(carSpeed)
	}

	if c.tryPark() {
		return
	}

	isStepValid := c.step >= 1
	c.angle = angle

	if !isOnTarget && isStepValid {
		c.x += speed
	} else {
		c.x = int(c.targetSlot.X)
		c.y = int(c.targetSlot.Y)
		switch c.targetSlot.Type {
		case "down":
			c.angle = 180
		case "left":
			c.angle = -90
		case "right":
			c.angle = 90
		case "up":
			c.angle = 0
		}
		c.internalStep = 5
	}
}

func (c *Car) tryPark() bool {
	nearDistance := 80.00
	if math.Abs(float64(c.x)-c.targetSlot.X) < nearDistance && math.Abs(float64(c.y)-c.targetSlot.Y) < nearDistance {
		c.x = int(c.targetSlot.X)
		c.y = int(c.targetSlot.Y)
		switch c.targetSlot.Type {
		case "down":
			c.angle = 180
		case "left":
			c.angle = -90
		case "right":
			c.angle = 90
		case "up":
			c.angle = 0
		}
		return true
	}
	return false
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

func (c *Car) Park() {
	c.step = 1
}
