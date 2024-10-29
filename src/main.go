package main

import "Parking_Simulator/src/core/render"

const (
	iconPath = "src/assets/icon.png"
	mapPath  = "src/assets/map/editable/map/main.tmx"
	appName  = "Parking Simulator"
)

func main() {

	resources :=
		render.NewResources(iconPath, mapPath, render.CarSpritePath{
			Blue:   "src/assets/sprites/car-blue.png",
			Green:  "src/assets/sprites/car-green.png",
			Orange: "src/assets/sprites/car-orange.png",
			Pink:   "src/assets/sprites/car-pink.png",
			Red:    "src/assets/sprites/car-red.png",
			Yellow: "src/assets/sprites/car-yellow.png",
		})

	renderEngine := render.NewRenderEngine(appName, resources.Icon, resources.MapInfo, resources.CarsSprites)
	renderEngine.Run()
}
