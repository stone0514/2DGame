package systems

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

//PlayerInit ...
func (ps *PlayerSystem) PlayerInit(player *Player) {

	PsPositionX := float32(0)
	PsPositionY := engo.WindowHeight() - CellHeight33*6

	//spaceComponent
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: PsPositionX, Y: PsPositionY},
		Width:    30,
		Height:   30,
	}

	//create spritesheet
	player.spritesheet = common.NewSpritesheetWithBorderFromFile(playerFile, 64, 66, 0, 0)

	//rendreComponent
	player.RenderComponent = common.RenderComponent{
		Drawable: player.spritesheet.Cell(PlayerSpriteSheetCell),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	player.RenderComponent.SetZIndex(5)

	//set component
	ps.playerEntity = player

	//initialize
	ps.playerEntity.playerPositionY = PsPositionY
	ps.playerEntity.LeftPositionX = PsPositionX + float32(ExtraSizeX)
	ps.playerEntity.RightPositionX = PsPositionX + CellWidth32 - float32(ExtraSizeX)
	ps.playerEntity.ifFalling = false
	ps.playerEntity.ifOnPipe = false
	ps.playerEntity.cameraMoveDistance = 0
	ps.playerEntity.topCount = 1 + MaxCount/2
	ps.playerEntity.bottomCount = 0
	ps.playerEntity.ifStart = false
	ifGameOver = false

	//add renderSystem
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}

	//move camera
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:        common.XAxis,
		Value:       engo.WindowWidth() / 2,
		Incremental: false,
	})
}
