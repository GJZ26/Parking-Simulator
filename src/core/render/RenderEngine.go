package render

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"image"
	"log"
)

const (
	appName  = "Parking Simulator"
	mapPath  = "src/assets/map/editable/map/main.tmx"
	iconPath = "src/assets/icon.png"
)

const (
	tilesAmount  = 13
	tilesSize    = 70
	screenWidth  = tilesSize * tilesAmount
	screenHeight = tilesSize * tilesAmount
)

type RenderEngine struct {
	background *ebiten.Image
}

func NewRenderEngine() *RenderEngine {
	renderEngine := &RenderEngine{}
	renderEngine.setWindow()
	renderEngine.loadMap(mapPath)
	return renderEngine
}

func (r *RenderEngine) loadMap(imagePath string) {
	fmt.Println("Loading map from ", imagePath)
	mapInfo, err := tiled.LoadFile(imagePath)
	if err != nil {
		fmt.Println("Error loading map: ", err)
		log.Fatal(err)
	}
	fmt.Println("Rendering map")

	rendered, err := render.NewRenderer(mapInfo)
	if err != nil {
		fmt.Println("Error rendering map: ", err)
		log.Fatal(err)
	}

	if err := rendered.RenderVisibleLayers(); err != nil {
		fmt.Println("Error rendering map: ", err)
	}

	r.background = ebiten.NewImageFromImage(rendered.Result)
}

func (r *RenderEngine) Draw(screen *ebiten.Image) {
	screen.DrawImage(r.background, &ebiten.DrawImageOptions{})
}

func (r *RenderEngine) Update() error {
	return nil
}

func (r *RenderEngine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (r *RenderEngine) setWindow() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(appName)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	icon, _, err := ebitenutil.NewImageFromFile(iconPath)
	if err != nil {
		fmt.Println("Error loading icon: ", err)
		log.Fatal(err)
	}

	ebiten.SetWindowIcon([]image.Image{icon})
}

func (r *RenderEngine) Run() {
	fmt.Println("Running render loop")
	if err := ebiten.RunGame(r); err != nil {
		log.Fatal(err)
	}
}
