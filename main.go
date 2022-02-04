package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type mySceane struct{}

func (*mySceane) Type() string { return "myGame" }

func (*mySceane) Preload() {
	engo.Files.Load("./assets/character/Adventurer/adventurer-idle-00.png")
	engo.Files.Load("./assets/tileSet/world/Asset_01/pixel_platform_01_tileset_final.png")
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
		Title:  "Hello World",
		Width:  1280,
		Height: 800,
	}
	engo.Run(opts, &mySceane{})
}
