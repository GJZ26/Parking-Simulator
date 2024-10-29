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
	Points      PointMap
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

type PointData struct {
	Id     uint32
	X      float64
	Y      float64
	Width  float64
	Height float64
	Type   string
}

type PointMap struct {
	Road        []PointData
	ParkingSlot []PointData
	ParkingRoad []PointData
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

	addPoint := func(points *[]PointData, object *tiled.Object) {
		*points = append(*points, PointData{
			Id:     object.ID,
			X:      object.X,
			Y:      object.Y,
			Width:  object.Width,
			Height: object.Height,
			Type:   object.Type,
		})
	}

	for _, layer := range tiledMap.ObjectGroups {
		switch layer.Name {
		case "slot":
			for _, object := range layer.Objects {
				addPoint(&r.Points.ParkingSlot, object)
			}
		case "path":
			for _, object := range layer.Objects {
				if object.Type == "default" {
					addPoint(&r.Points.Road, object)
				} else {
					addPoint(&r.Points.ParkingRoad, object)
				}
			}
		}
	}
}
