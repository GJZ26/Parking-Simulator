package routines

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager/types"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type EntityManager struct {
	cars       []*entity.Car
	carSprites *types.CarSprites
	points     types.PointMap
}

func NewEntityManager(
	points types.PointMap,
	sprites types.CarSprites) *EntityManager {
	return &EntityManager{
		cars:       make([]*entity.Car, 0),
		points:     points,
		carSprites: &sprites,
	}
}

func (e *EntityManager) invokeCar(amount int) {
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

	for i := 0; i < amount; i++ {
		e.cars = append(e.cars, entity.NewCar(rng.Intn(1000), rng.Intn(1000), getRandomSprite()))
	}
}

func (e *EntityManager) Run(renderChannel chan []*entity.Car) {
	if len(e.cars) <= 0 {
		fmt.Println("Invoking!")
		e.invokeCar(10)
	}
	fmt.Println("a")
	renderChannel <- e.cars
}
