package render

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"image"
	"log"
)

const (
	tilesAmount  = 13
	tilesSize    = 70
	screenWidth  = tilesSize * tilesAmount
	screenHeight = tilesSize * tilesAmount
)

type Engine struct {
	background *ebiten.Image
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
