package manager

import (
	"Parking_Simulator/src/core/manager/types"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"log"
	"reflect"
)

type Resources struct {
	Icon        *ebiten.Image
	MapInfo     *tiled.Map
	CarsSprites types.CarSprites
	Points      types.PointMap
}

func NewResources(iconPath string, mapPath string, carPaths types.CarSpritePath) *Resources {
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
	resources.getPoints(mapInfo)

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

func (r *Resources) getPoints(tiledMap *tiled.Map) {
	fmt.Println("Reading object points from Map")

	addPoint := func(points *[]types.PointData, object *tiled.Object) {
		if points != nil {
			*points = append(*points, types.PointData{
				Id:     object.ID,
				X:      object.X,
				Y:      object.Y,
				Width:  object.Width,
				Height: object.Height,
				Type:   object.Type,
			})
		} else {
			fmt.Println("Warning: Points array is nil")
		}
	}

	for _, layer := range tiledMap.ObjectGroups {
		for _, object := range layer.Objects {
			switch layer.Name {
			case "default-path":
				addPoint(&r.Points.Road, object)
			case "parking-path":
				addPoint(&r.Points.ParkingRoad, object)
			case "slot-right":
				addPoint(&r.Points.RightParkingSlot, object)
			case "slot-left":
				addPoint(&r.Points.LeftParkingSlot, object)
			default:
				fmt.Println("Warning: Unknown object")
			}
		}
	}
}
