package routines

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager/types/geo"
	"Parking_Simulator/src/core/manager/types/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

const (
	maxCar        = 100
	maxCarInQueue = 3
)

type EntityManager struct {
	carCounter          int
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
		idleCars:      make([]*entity.Car, 0),
		entranceQueue: make([]*entity.Car, 0),
		exitQueue:     make([]*entity.Car, 0),
		points:        points,
		carSprites:    &sprites,
	}
}

func (e *EntityManager) invokeCar(slot geo.SlotInfo, queue *[]*entity.Car) {
	println("[Entity Manager]: Invoking car")
	if e.carCounter >= maxCar {
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
		default:
		}

		renderChannel <- resources.RenderData{
			Counter: maxCar - e.carCounter,
			Cars:    e.tick(freeSlotChannel, exitControl),
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (e *EntityManager) tick(freeSlotChannel chan []uint32, exitControl chan bool) []*entity.Car {
	var freeSlots = make([]uint32, 0)
	// Update entrance cars
	for _, car := range e.entranceQueue {
		car.Update()
	}

	// Update idle cars
	for index := 0; index < len(e.idleCars); {
		car := e.idleCars[index]
		stillParking, needDead := car.Update()
		if car.IsInOtherSide() {
			select {
			case exitControl <- true:
			default:
				println("- - - [Entity Manager]: ERROR, COLLISION!! - - -")
			}
		}
		if !stillParking {
			if len(e.exitQueue) < maxCarInQueue {
				car.ToExitQueue()
				e.exitQueue = append(e.exitQueue, car)
				e.idleCars = append(e.idleCars[:index], e.idleCars[index+1:]...)

				for i, waitedCar := range e.exitQueue {
					waitedCar.SetQueuePosition(i)
				}

				freeSlots = append(freeSlots, car.GetSlotID())
			}
		}
		if needDead {
			e.idleCars = append(e.idleCars[:index], e.idleCars[index+1:]...)
		} else {
			index++
		}
	}

	for _, car := range e.exitQueue {
		car.Update()
	}

	select {
	case freeSlotChannel <- freeSlots:
	default:
	}

	return append(e.entranceQueue, append(e.idleCars, e.exitQueue...)...)
}

func (e *EntityManager) newSlotAvailable(info geo.SlotInfo) {
	if len(e.entranceQueue) < maxCarInQueue {
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
			e.processExit()
		}
		return
	}

	if len(e.entranceQueue) > 0 {
		e.processEntrance()
		return
	}

	if len(e.exitQueue) > 0 {
		e.processExit()
	}
}

func (e *EntityManager) processEntrance() {
	car := e.entranceQueue[0]
	car.Park()

	e.entranceQueue = e.entranceQueue[1:]
	e.idleCars = append(e.idleCars, car)

	if len(e.slotsAvailable) > 0 {
		e.invokeCar(e.slotsAvailable[0], &e.entranceQueue)
		e.slotsAvailable = e.slotsAvailable[1:]
	}

	for i, c := range e.entranceQueue {
		c.SetQueuePosition(i)
	}
}

func (e *EntityManager) processExit() {
	car := e.exitQueue[0]
	car.Exit()

	e.exitQueue = e.exitQueue[1:]
	e.idleCars = append(e.idleCars, car)

	for i, c := range e.exitQueue {
		c.SetQueuePosition(i)
	}
}
