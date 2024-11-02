package render

import (
	"Parking_Simulator/src/core/entity"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"image"
	"log"
	"math"
)

const (
	tilesAmount  = 13
	tilesSize    = 70
	screenWidth  = tilesSize * tilesAmount
	screenHeight = tilesSize * tilesAmount
)

type Engine struct {
	background    *ebiten.Image
	renderChannel chan []*entity.Car
	cache         []*entity.Car
}

func NewRenderEngine(title string, icon *ebiten.Image, tileMap *tiled.Map, renderChannel chan []*entity.Car) *Engine {
	renderEngine := &Engine{
		renderChannel: renderChannel,
	}
	renderEngine.setWindow(title, icon)
	renderEngine.loadMap(tileMap)
	return renderEngine
}

func (r *Engine) loadMap(tileMap *tiled.Map) {
	rendered, err := render.NewRenderer(tileMap)
	if err != nil {
		fmt.Println("Error rendering map: ", err)
		log.Fatal(err)
	}

	if err := rendered.RenderVisibleLayers(); err != nil {
		fmt.Println("Error rendering map: ", err)
	}

	r.background = ebiten.NewImageFromImage(rendered.Result)
}

func (r *Engine) Draw(screen *ebiten.Image) {
	screen.DrawImage(r.background, &ebiten.DrawImageOptions{})
	select {
	case current := <-r.renderChannel:
		// Render from channel
		r.DrawCars(screen, current)
		r.cache = current
	default:
		// Render from caché
		r.DrawCars(screen, r.cache)
	}
}

func (r *Engine) DrawCars(screen *ebiten.Image, cars []*entity.Car) {
	for _, car := range cars {
		r.DrawCar(screen, car)
	}
}

func (r *Engine) DrawCar(screen *ebiten.Image, car *entity.Car) {

	s := car.GetSprite().Bounds().Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)

	op.GeoM.Rotate(float64(car.GetAngle()) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(car.GetX()), float64(car.GetY()))

	screen.DrawImage(car.GetSprite(), op)
}

func (r *Engine) Update() error {
	return nil
}

func (r *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (r *Engine) setWindow(title string, icon *ebiten.Image) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{icon})
}

func (r *Engine) Run() {
	fmt.Println("Running render loop")
	if err := ebiten.RunGame(r); err != nil {
		log.Fatal(err)
	}
}
