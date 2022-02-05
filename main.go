package main

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type mySceane struct{}

func (*mySceane) Type() string { return "myGame" }

func (*mySceane) Preload() {
	engo.Files.Load("./assets/character/Adventurer/adventurer-Sheet.png")
	engo.Files.Load("./assets/character/Adventurer/adventurer-idle-00.png")
	engo.Files.Load("./assets/tileSet/world/Asset_01/pixel_platform_01_tileset_final.png")
	common.SetBackground(color.RGBA{210, 180, 140, 1})
}

func (*mySceane) Setup(u engo.Updater) {
	engo.Input.RegisterButton("Right", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("Left", engo.KeyA, engo.KeyArrowRight)
	engo.Input.RegisterButton("Jump", engo.KeySpace)

	world, _ := u.(*ecs.World)

	world.AddSystem(&common.RenderSystem{})
	/*
		world.AddSystem(&systems.TileSystem{})
		world.AddSystem(&systems.PlayerSystem{})
	*/

}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          1280,
		Height:         800,
		StandardInputs: true,
		NotResizable:   true,
	}
	fmt.Println("START AdventureGame")
	engo.Run(opts, &mySceane{})
}

func (*mySceane) Exit() {
	engo.Exit()
}
