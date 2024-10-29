package render

import (
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
	background *ebiten.Image
	carAssets  CarSprites
}

func NewRenderEngine(title string, icon *ebiten.Image, tileMap *tiled.Map, carSprites CarSprites) *Engine {
	renderEngine := &Engine{}
	renderEngine.setWindow(title, icon)
	renderEngine.loadMap(tileMap)
	renderEngine.carAssets = carSprites
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
	s := r.carAssets.Blue.Bounds().Size()
	op := &ebiten.DrawImageOptions{}
	// Move the image's center to the screen's upper-left corner.
	// This is a preparation for rotating. When geometry matrices are applied,
	// the origin point is the upper-left corner.
	op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)

	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	op.GeoM.Rotate(float64(0) * 2 * math.Pi / 360)

	// Move the image to the screen's center.
	op.GeoM.Translate(0, 0)

	screen.DrawImage(r.background, &ebiten.DrawImageOptions{})
	screen.DrawImage(r.carAssets.Blue, op)
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
