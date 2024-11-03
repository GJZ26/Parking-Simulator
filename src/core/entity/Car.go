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
	exitStep      uint8
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

func (c *Car) Update() (bool, bool) {
	switch c.internalStep {
	case 0:
		c.posInEntranceQueue()
	case 1:
		c.goToEntrance()
	case 2:
		c.goToCross()
	case 3, 4, 5:
		c.findWayToSlot()
	case 6:
		c.searchExitCross()
	case 7:
		c.exit()
	default:
		return c.IsActive(), c.internalStep == 8
	}

	return c.IsActive(), c.internalStep == 8
}

func (c *Car) searchExitCross() {
	if c.step <= 1 {
		return
	}
	if c.exitStep == 0 {
		c.exitStep = 1
	}
	switch c.exitStep {
	case 1:
		c.searchNearestPath()
	case 2:
		c.goToBottomCorner()
	case 3:
		c.goToUpperCorner()
	case 4:
		c.goToExit()
	case 5:
		c.formToQueue()

	}
}

func (c *Car) exit() {
	switch c.exitStep {
	case 0:
		c.goOutSide()
	case 1:
		c.runAway()
	}
}

func (c *Car) goOutSide() {
	target := c.defaultRoad[1].Y
	isOnPosition := c.y >= int(target)
	if !isOnPosition {
		c.y += carSpeed
	} else {
		c.angle = -90
		c.exitStep++
	}
}

func (c *Car) runAway() {
	target := c.defaultRoad[2].X
	isOnPosition := c.x <= int(target)
	if !isOnPosition {
		c.x -= carSpeed
	} else {
		c.internalStep++
	}
}

func (c *Car) searchNearestPath() {
	if c.targetSlot.Side == geo.LeftSide {
		if c.targetSlot.Type == "up" || c.targetSlot.Type == "down" {
			if math.Abs(c.parkingRoad[1].Y-float64(c.y)) < 78.00 {
				c.exitStep = 2
				c.y = int(c.parkingRoad[1].Y)
				c.angle = -90
				return
			}
			if math.Abs(c.parkingRoad[2].Y-float64(c.y)) < 78.00 {
				c.exitStep = 4
				c.y = int(c.parkingRoad[2].Y)
				c.angle = 90
				return
			}
		} else {
			c.x = int(c.parkingRoad[2].X)
			c.exitStep = 3
			c.angle = 0
			return
		}
	} else {
		if c.targetSlot.Type == "up" || c.targetSlot.Type == "down" {
			if math.Abs(c.parkingRoad[3].Y-float64(c.y)) < 78.00 {
				c.exitStep = 2
				c.y = int(c.parkingRoad[3].Y)
				c.angle = 90
				return
			}
			if math.Abs(c.parkingRoad[4].Y-float64(c.y)) < 78.00 {
				c.exitStep = 4
				c.y = int(c.parkingRoad[4].Y)
				c.angle = -90
				return
			}
		} else {
			c.x = int(c.parkingRoad[4].X)
			c.exitStep = 3
			c.angle = 0
			return
		}
	}
}

func (c *Car) goToBottomCorner() {
	var targetPosition int
	var isOnTarget bool
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		targetPosition = int(c.parkingRoad[1].X)
		isOnTarget = c.x <= targetPosition
		speed = -(carSpeed)
	} else {
		targetPosition = int(c.parkingRoad[3].X)
		isOnTarget = c.x >= targetPosition
		speed = carSpeed
	}

	if !isOnTarget {
		c.x += speed
	} else {
		c.angle = 0
		c.exitStep++
	}
}

func (c *Car) goToUpperCorner() {
	var targetPosition int
	var isOnTarget bool
	var speed int

	targetPosition = int(c.parkingRoad[2].Y)
	isOnTarget = c.y <= targetPosition
	speed = -(carSpeed)

	if !isOnTarget {
		c.y += speed
	} else {
		if c.targetSlot.Side == geo.LeftSide {
			c.angle = 90
		} else {
			c.angle = -90
		}

		c.exitStep++
	}

}

func (c *Car) goToExit() {
	targetPosition := int(c.parkingRoad[5].X)
	var isOnTarget bool
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		isOnTarget = c.x >= targetPosition
		speed = carSpeed
	} else {
		isOnTarget = c.x <= targetPosition
		speed = -carSpeed
	}

	if !isOnTarget {
		c.x += speed
	} else {
		c.angle = 180
		c.exitStep++
	}
}

func (c *Car) formToQueue() {
	targetPosition := int(c.parkingRoad[0].Y - 90 - (70 * float64(c.queuePosition)))
	isOnTarget := c.y >= targetPosition

	if !isOnTarget {
		c.y += carSpeed
	} else {
		c.internalStep++
		c.exitStep = 0
	}
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
	if !isStepValid {
		return
	}

	if isBeyondTarget {
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

	if !isStepValid {
		return
	}

	if isUnderTarget {
		c.y -= carSpeed
	} else {
		c.internalStep = 3
	}
}

// TODO: Clean up
func (c *Car) findWayToSlot() {
	var targetPosition float64
	var isOnTarget bool
	var angle float64
	var speed int

	if c.targetSlot.Side == geo.LeftSide {
		if c.internalStep == 3 {
			targetPosition = c.parkingRoad[1].X
			isOnTarget = c.x <= int(targetPosition)
			angle = -90
			speed = -(carSpeed)
		}
		if c.internalStep == 4 {
			targetPosition = c.parkingRoad[2].Y
			isOnTarget = c.y <= int(targetPosition)
			angle = 0
			speed = -(carSpeed)
		}
		if c.internalStep == 5 {
			targetPosition = c.parkingRoad[5].X
			isOnTarget = c.x >= int(targetPosition)
			angle = 90
			speed = carSpeed
		}
	} else {
		if c.internalStep == 3 {
			targetPosition = c.parkingRoad[3].X
			isOnTarget = c.x >= int(targetPosition)
			angle = 90
			speed = carSpeed
		}
		if c.internalStep == 4 {
			targetPosition = c.parkingRoad[4].Y
			isOnTarget = c.y <= int(targetPosition)
			angle = 0
			speed = -(carSpeed)
		}
		if c.internalStep == 5 {
			targetPosition = c.parkingRoad[5].X
			isOnTarget = c.x <= int(targetPosition)
			angle = -90
			speed = -(carSpeed)
		}
	}

	if c.tryPark(false) {
		c.internalStep = 6
		return
	}

	isStepValid := c.step >= 1
	c.angle = angle

	if !isOnTarget && isStepValid {
		if c.internalStep == 3 {
			c.x += speed
		}
		if c.internalStep == 4 {
			c.y += speed
		}
		if c.internalStep == 5 {
			c.x += speed
		}
	} else {
		if c.internalStep == 5 {
			c.tryPark(true)
			c.internalStep = 6
			return
		}
		c.internalStep++
	}
}

func (c *Car) tryPark(force bool) bool {
	nearDistance := 78.00
	if force || math.Abs(float64(c.x)-c.targetSlot.X) < nearDistance && math.Abs(float64(c.y)-c.targetSlot.Y) < nearDistance {
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
		c.StartTimer()
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

func (c *Car) ToExitQueue() {
	c.step = 2
}

func (c *Car) Exit() {
	c.step = 3
}
