package main

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager"
	"Parking_Simulator/src/core/manager/routines"
	"Parking_Simulator/src/core/manager/types"
	render2 "Parking_Simulator/src/core/render"
)

const (
	iconPath   = "src/assets/icon.png"
	mapPath    = "src/assets/map/editable/map/main.tmx"
	appName    = "Parking Simulator"
	bluePath   = "src/assets/sprites/car-blue.png"
	greenPath  = "src/assets/sprites/car-green.png"
	orangePath = "src/assets/sprites/car-orange.png"
	pinkPath   = "src/assets/sprites/car-pink.png"
	redPath    = "src/assets/sprites/car-red.png"
	yellowPath = "src/assets/sprites/car-yellow.png"
)

func main() {
	// Shared resources.
	var resources = manager.NewResources(iconPath, mapPath, types.CarSpritePath{
		Blue:   bluePath,
		Green:  greenPath,
		Orange: orangePath,
		Pink:   pinkPath,
		Red:    redPath,
		Yellow: yellowPath,
	})

	// Starting Entity Manager
	var entityManager = routines.NewEntityManager(resources.Points.ParkingRoad, resources.Points.ParkingSlot, resources.Points.Road, resources.CarsSprites)

	// Channels
	renderChannel := make(chan []*entity.Car)

	go entityManager.Run(renderChannel)

	// Starting Render Engine (Ebiten)
	var renderEngine = render2.NewRenderEngine(appName, resources.Icon, resources.MapInfo, renderChannel)
	renderEngine.Run()
}
