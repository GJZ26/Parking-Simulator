package main

import (
	"Parking_Simulator/src/core/manager"
	"Parking_Simulator/src/core/manager/routines"
	"Parking_Simulator/src/core/manager/types/geo"
	resourcesType "Parking_Simulator/src/core/manager/types/resources"
	render2 "Parking_Simulator/src/core/render"
	"fmt"
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
	fmt.Println("Starting", appName)
	// Shared resources.
	var resources = manager.NewResources(iconPath, mapPath, resourcesType.CarSpritePath{
		Blue:   bluePath,
		Green:  greenPath,
		Orange: orangePath,
		Pink:   pinkPath,
		Red:    redPath,
		Yellow: yellowPath,
	})

	// Starting Entity Manager
	var entityManager = routines.NewEntityManager(resources.Points, resources.CarsSprites)
	var slotManager = routines.NewSlotManager(resources.Points)
	var renderEngine = render2.NewRenderEngine(appName, resources.Icon, resources.MapInfo)

	// Channels
	slotChannel := make(chan geo.SlotInfo, 1)
	renderChannel := make(chan resourcesType.RenderData, 1)
	freeSlotChannel := make(chan []uint32, 1)
	exitChannel := make(chan bool)
	entranceChannel := make(chan bool)

	// GoRoutines
	go routines.ExitManager(exitChannel, entranceChannel)
	go entityManager.Run(renderChannel, slotChannel, freeSlotChannel, exitChannel, entranceChannel)
	go slotManager.Run(slotChannel, freeSlotChannel)
	go renderEngine.UpdateCache(renderChannel)

	// Render loop :D
	renderEngine.Run()
}
