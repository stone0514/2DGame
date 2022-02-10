package main

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/stone0514/2DGame/systems"
)

type mySceane struct{}

func (*mySceane) Type() string { return "myGame" }

func (*mySceane) Preload() {
	engo.Files.Load("./characters/robot.png")              //48x48
	engo.Files.Load("./characters/robot2.png")             //64x64
	engo.Files.Load("./tileSet/world/tileSpritesheet.png") //33x33
	engo.Files.Load("./tileSet/world/spritesheet_sample.png")
	engo.Files.Load("./tileSet/world/Background/clouds.png")
	common.SetBackground(color.RGBA{210, 180, 140, 1})
}

func (*mySceane) Setup(u engo.Updater) {
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("Jump", engo.KeySpace)

	world, _ := u.(*ecs.World)

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.TileSystem{})
	world.AddSystem(&systems.PlayerSystem{})

}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          800,
		Height:         600,
		StandardInputs: true,
		NotResizable:   true,
	}
	fmt.Println("START AdventureGame")
	engo.Run(opts, &mySceane{})
}

func (*mySceane) Exit() {
	engo.Exit()
}
