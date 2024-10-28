package render

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"log"
	"reflect"
)

type ResourcesRequest struct {
	iconPath string
}

type Resources struct {
	Icon        *ebiten.Image
	MapInfo     *tiled.Map
	CarsSprites CarSprites
}

type CarSprites struct {
	Blue   *ebiten.Image
	Green  *ebiten.Image
	Orange *ebiten.Image
	Pink   *ebiten.Image
	Red    *ebiten.Image
	Yellow *ebiten.Image
}

type CarSpritePath struct {
	Blue   string
	Green  string
	Orange string
	Pink   string
	Red    string
	Yellow string
}

func NewResources(iconPath string, mapPath string, carPaths CarSpritePath) *Resources {
	fmt.Println("Importing resources")
	resources := Resources{}

	icon, _, err := ebitenutil.NewImageFromFile(iconPath)
	if err != nil {
		fmt.Println("Error loading icon")
		log.Fatal(err)
	}
	resources.Icon = icon

	mapInfo, err := tiled.LoadFile(mapPath)
	if err != nil {
		fmt.Println("Error loading map")
		log.Fatal(err)
	}
	resources.MapInfo = mapInfo

	carPathsValue := reflect.ValueOf(carPaths)
	carSpritesValue := reflect.ValueOf(&resources.CarsSprites).Elem()

	for i := 0; i < carPathsValue.NumField(); i++ {
		colorName := carPathsValue.Type().Field(i).Name
		colorPath := carPathsValue.Field(i).String()

		carImage, _, err := ebitenutil.NewImageFromFile(colorPath)
		if err != nil {
			fmt.Printf("Error loading car sprite for %s\n", colorName)
			log.Fatal(err)
		}

		carSpritesValue.FieldByName(colorName).Set(reflect.ValueOf(carImage))
		fmt.Printf("Loaded car sprite for %s: %s\n", colorName, colorPath)
	}

	return &resources
}
