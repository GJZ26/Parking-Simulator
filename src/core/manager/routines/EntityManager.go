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
	exitControl         bool
	entranceQueue       []*entity.Car
	exitQueue           []*entity.Car
	idleCars            []*entity.Car
	carSprites          *resources.CarSprites
	points              geo.PointMap
	slotsAvailable      []geo.SlotInfo
}

func NewEntityManager(
	points geo.PointMap,
	sprites resources.CarSprites) *EntityManager {
	return &EntityManager{
		carCounter:    0,
		maxCars:       100,
		idleCars:      make([]*entity.Car, 0),
		entranceQueue: make([]*entity.Car, 0),
		exitQueue:     make([]*entity.Car, 0),
		points:        points,
		carSprites:    &sprites,
	}
}

func (e *EntityManager) invokeCar(slot geo.SlotInfo, queue *[]*entity.Car) {
	println("[Entity Manager]: Invoking car")
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
	car := entity.NewCar(slot, e.points.Road, e.points.ParkingRoad, getRandomSprite(), len(*queue))
	*queue = append(*queue, car)
	e.carCounter++
}

func (e *EntityManager) Run(renderChannel chan resources.RenderData, slotChannel chan geo.SlotInfo, freeSlotChannel chan []uint32, exitControl chan bool, entranceControl chan bool) {
	for {
		select {
		case slot := <-slotChannel:
			e.newSlotAvailable(slot)
		case <-entranceControl:
			e.entryOrExit()
			exitControl <- true
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
	activeCars := make([]*entity.Car, 0, len(e.entranceQueue))
	var res = make([]uint32, 0)

	for _, car := range e.idleCars {
		stillParking, needDead := car.Update()
		if !needDead {
			activeCars = append(activeCars, car)
		}
		if !stillParking {
			car.Exit()
		}
	}
	e.idleCars = activeCars

	activeCars = make([]*entity.Car, 0, len(e.entranceQueue))
	for _, car := range e.entranceQueue {
		needExit, _ := car.Update()
		if needExit {
			activeCars = append(activeCars, car)
		} else {
			res = append(res, car.GetSlotID())
		}
	}
	e.entranceQueue = activeCars

	select {
	case freeSlotChannel <- res:
	default:
		println("[EntityManager]: Canal de freeSlotChannel no disponible, datos no enviados")
	}

	return append(e.entranceQueue, append(e.idleCars, e.exitQueue...)...)
}

func (e *EntityManager) newSlotAvailable(info geo.SlotInfo) {
	if len(e.entranceQueue) < 3 {
		if len(e.slotsAvailable) > 0 {
			info = e.slotsAvailable[len(e.slotsAvailable)-1]
			e.slotsAvailable = e.slotsAvailable[:len(e.slotsAvailable)-1]
		}
		e.invokeCar(info, &e.entranceQueue)
	} else {
		e.slotsAvailable = append(e.slotsAvailable, info)
	}
}

func (e *EntityManager) entryOrExit() {
	if len(e.entranceQueue) > 0 && len(e.exitQueue) > 0 {
		if rand.Intn(2) == 0 {
			e.processEntrance()
		} else {
			//e.processExit()
		}
		return
	}

	if len(e.entranceQueue) > 0 {
		e.processEntrance()
		return
	}

	if len(e.exitQueue) > 0 {
		// e.processExit()
	}
}

func (e *EntityManager) processEntrance() {
	car := e.entranceQueue[0]
	car.Park()

	e.entranceQueue = e.entranceQueue[1:]
	e.idleCars = append(e.idleCars, car)

	for i, c := range e.entranceQueue {
		c.SetQueuePosition(i)
	}
}
