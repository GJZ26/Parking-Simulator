package routines

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager/types/geo"
	"Parking_Simulator/src/core/manager/types/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type EntityManager struct {
	carCounter          int
	maxCars             int
	carsInEntranceQueue int
	carsInExitQueue     int
	cars                []*entity.Car
	carSprites          *resources.CarSprites
	points              geo.PointMap
	slotsAvailable      []geo.SlotInfo
}

func NewEntityManager(
	points geo.PointMap,
	sprites resources.CarSprites) *EntityManager {
	return &EntityManager{
		carsInEntranceQueue: 0,
		carsInExitQueue:     0,
		carCounter:          0,
		maxCars:             100,
		cars:                make([]*entity.Car, 0),
		points:              points,
		carSprites:          &sprites,
	}
}

func (e *EntityManager) invokeCar(slot geo.SlotInfo, x int, y int) {
	if e.carCounter >= e.maxCars {
		return
	}
	rng := rand.New(rand.NewSource(rand.Int63()))

	getRandomSprite := func() *ebiten.Image {
		switch rng.Intn(6) {
		case 0:
			return e.carSprites.Blue
		case 1:
			return e.carSprites.Green
		case 2:
			return e.carSprites.Orange
		case 3:
			return e.carSprites.Pink
		case 4:
			return e.carSprites.Red
		default:
			return e.carSprites.Yellow
		}
	}
	println("[Entity Manager]: Invoking new Car")
	e.cars = append(e.cars, entity.NewCar(slot, e.points.Road, e.points.ParkingRoad, getRandomSprite()))
	e.carCounter++
}

func (e *EntityManager) Run(renderChannel chan resources.RenderData, slotChannel chan geo.SlotInfo, freeSlotChannel chan []uint32) {
	for {
		select {
		case slot := <-slotChannel:
			e.newSlotAvailable(slot)
		default:
		}

		renderChannel <- resources.RenderData{
			Counter: e.maxCars - e.carCounter,
			Cars:    e.Update(freeSlotChannel),
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (e *EntityManager) Update(freeSlotChannel chan []uint32) []*entity.Car {
	activeCars := make([]*entity.Car, 0, len(e.cars))
	var res = make([]uint32, 0)

	for _, car := range e.cars {
		if car.Update() {
			activeCars = append(activeCars, car)
		} else {
			res = append(res, car.GetSlotID())
		}
	}
	e.cars = activeCars

	select {
	case freeSlotChannel <- res:
	default:
		println("[EntityManager]: Canal de freeSlotChannel no disponible, datos no enviados")
	}

	return e.cars
}

func (e *EntityManager) newSlotAvailable(info geo.SlotInfo) {
	if e.carCounter < 3 {
		if len(e.slotsAvailable) > 0 {
			info = e.slotsAvailable[len(e.slotsAvailable)-1]
			e.slotsAvailable = e.slotsAvailable[:len(e.slotsAvailable)-1]
		}
		e.invokeCar(info, int(info.X), int(info.Y))
		e.carCounter++
	} else {
		e.slotsAvailable = append(e.slotsAvailable, info)
	}
}
