package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/stone0514/2DGame/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type mySceane struct{}

func (*mySceane) Type() string { return "myGame" }

//preload pngFiles
func (*mySceane) Preload() {
	engo.Files.Load("./characters/robot.png")              //w:64xh:66
	engo.Files.Load("./tileSet/world/tileSpritesheet.png") //w:33xh:33
	engo.Files.Load("./tileSet/world/spritesheet_sample.png")
	engo.Files.Load("./tileSet/world/Background/clouds.png") //w:98xh:60
	/*
		engo.Files.Load("./vehicle/vehicle-1.png")
		engo.Files.Load("./vehicle/vehicle-2.png")
		engo.Files.Load("./vehicle/vehicle-3.png")
	*/
	common.SetBackground(color.RGBA{106, 90, 205, 1})
}

//registerButton,add systems to scene
func (*mySceane) Setup(u engo.Updater) {
	//key config
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("Jump", engo.KeySpace)
	engo.Input.RegisterButton("Enter", engo.KeyEnter)
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))

	//world config
	world, _ := u.(*ecs.World)

	//add system
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.TileSystem{})
	world.AddSystem(&systems.PlayerSystem{})
	world.AddSystem(&systems.EnermySystem{})
	world.AddSystem(&systems.HUDTextSystem{})

}

//main
func main() {
	opts := engo.RunOptions{
		Title:          "Mr.Robot",
		Width:          800,
		Height:         500,
		StandardInputs: true,
		NotResizable:   true,
	}
	fmt.Println("RUN!")
	engo.Run(opts, &mySceane{})
}

func (*mySceane) Exit() {
	engo.Exit()
}
