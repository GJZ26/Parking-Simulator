package routines

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager/types"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type EntityManager struct {
	carCounter int
	maxCars    int
	cars       []*entity.Car
	carSprites *types.CarSprites
	points     types.PointMap
}

func NewEntityManager(
	points types.PointMap,
	sprites types.CarSprites) *EntityManager {
	return &EntityManager{
		carCounter: 0,
		maxCars:    100,
		cars:       make([]*entity.Car, 0),
		points:     points,
		carSprites: &sprites,
	}
}

func (e *EntityManager) invokeCar(slotId uint32, x int, y int) {
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
	e.cars = append(e.cars, entity.NewCar(x, y, slotId, getRandomSprite()))
	e.carCounter++
}

func (e *EntityManager) Run(renderChannel chan types.RenderData, slotChannel chan types.SlotInfo, freeSlotChannel chan []uint32) {
	for {
		select {
		case slot := <-slotChannel:
			println("[Entity Manager]: Received free slot")
			e.invokeCar(slot.Id, int(slot.X), int(slot.Y))
		default:
		}
		println("[Entity Manager]: Send Render Data")
		renderChannel <- types.RenderData{
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
	freeSlotChannel <- res
	return e.cars
}
