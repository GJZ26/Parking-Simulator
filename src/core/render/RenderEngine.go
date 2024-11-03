package render

import (
	"Parking_Simulator/src/core/entity"
	"Parking_Simulator/src/core/manager/types"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	background *ebiten.Image
	cache      types.RenderData
}

func NewRenderEngine(title string, icon *ebiten.Image, tileMap *tiled.Map) *Engine {
	renderEngine := &Engine{}
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
	r.DrawCars(screen, r.cache)
}

func (r *Engine) DrawCars(screen *ebiten.Image, cars types.RenderData) {
	for _, car := range cars.Cars {
		r.DrawCar(screen, car)
	}

	text := fmt.Sprintf("Carros por estacionar: %d", cars.Counter)
	ebitenutil.DebugPrintAt(screen, text, 10, 10)
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

func (r *Engine) UpdateCache(renderChannel chan types.RenderData) {
	for {
		select {
		case current := <-renderChannel:
			println("[Render Engine]: Rendering from Channel")
			r.cache = current
		default:
			continue
		}
	}
}
