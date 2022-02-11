package systems

import (
	//"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	MoveDistance          = 4
	JumpHeight            = 4
	MaxCount              = 40
	PlayerSpriteSheetCell = 0
	ExtraSizeX            = 0
)

var playerFile = "./characters/robot.png"

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	playerPositionY    float32
	LeftPositionX      float32
	RightPositionX     float32
	cameraMoveDistance int
	spritesheet        *common.Spritesheet
	useCell            int
	jumpCount          int
	jumpCount2Step     int
	topCount           int
	bottomCount        int
	ifJumping          bool
	ifFalling          bool
}

type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
}

func (ps *PlayerSystem) New(w *ecs.World) {
	ps.world = w
	player := Player{BasicEntity: ecs.NewBasic()}

	PsPositionX := float32(0)
	PsPositionY := engo.WindowHeight() - CellHeight33*6

	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: PsPositionX, Y: PsPositionY},
		Width:    30,
		Height:   30,
	}

	player.spritesheet = common.NewSpritesheetWithBorderFromFile(playerFile, 64, 66, 0, 0)
	player.RenderComponent = common.RenderComponent{
		Drawable: player.spritesheet.Cell(PlayerSpriteSheetCell),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	player.RenderComponent.SetZIndex(0)

	ps.playerEntity = &player

	ps.playerEntity.playerPositionY = PsPositionY
	ps.playerEntity.LeftPositionX = PsPositionX + float32(ExtraSizeX)
	ps.playerEntity.RightPositionX = PsPositionX + CellWidth32 - float32(ExtraSizeX)
	ps.playerEntity.ifFalling = false
	ps.playerEntity.cameraMoveDistance = 0
	ps.playerEntity.topCount = 1 + MaxCount/2
	ps.playerEntity.bottomCount = 0

	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 3200, Y: 300},
	}
}

func (ps *PlayerSystem) Update(dt float32) {

	if ps.playerEntity.SpaceComponent.Position.Y == ps.playerEntity.playerPositionY {

		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell)
	}

	if int(ps.playerEntity.LeftPositionX) >= (TileNum-GoalTileNum+2)*CellWidth33 {
		ps.Remove(ps.playerEntity.BasicEntity)
	}

	if engo.Input.Button("MoveRight").Down() {
		if int(ps.playerEntity.SpaceComponent.Position.X) < ps.playerEntity.cameraMoveDistance+int(engo.WindowWidth())/2 {
			ps.playerEntity.SpaceComponent.Position.X += MoveDistance
			ps.playerEntity.LeftPositionX += MoveDistance
			ps.playerEntity.RightPositionX += MoveDistance
		} else {
			if int(ps.playerEntity.SpaceComponent.Position.X) < int(engo.WindowWidth())-CellWidth32 {
				ps.playerEntity.SpaceComponent.Position.X += MoveDistance
				ps.playerEntity.LeftPositionX += MoveDistance
				ps.playerEntity.RightPositionX += MoveDistance
			}
			if int(ps.playerEntity.SpaceComponent.Position.X) < TileNum*CellWidth33-int(engo.WindowWidth())/2 {
				engo.Mailbox.Dispatch(common.CameraMessage{
					Axis:        common.XAxis,
					Value:       MoveDistance,
					Incremental: true,
				})
			}
			ps.playerEntity.cameraMoveDistance += MoveDistance
		}
		/*
				for i := 0; i <= 13; i++ {
					if ps.playerEntity.jumpCount == 0 {
						ps.playerEntity.useCell++
						ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell + ps.playerEntity.useCell)
						if ps.playerEntity.useCell == 14 {
							for j := 0; j <= 13; j-- {
								if ps.playerEntity.useCell == 0 {
									break
								}
								ps.playerEntity.useCell--
								ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell + ps.playerEntity.useCell)
							}
						}
					}
				}
			}
		*/
		if ps.playerEntity.jumpCount == 0 {
			switch ps.playerEntity.useCell {
			case 0:
				ps.playerEntity.useCell = 3
			case 1:
				ps.playerEntity.useCell = 4
			case 2:
				ps.playerEntity.useCell = 5
			case 3:
				ps.playerEntity.useCell = 6
			case 4:
				ps.playerEntity.useCell = 7
			case 5:
				ps.playerEntity.useCell = 8
			case 6:
				ps.playerEntity.useCell = 9
			case 7:
				ps.playerEntity.useCell = 10
			case 8:
				ps.playerEntity.useCell = 12
			case 9:
				ps.playerEntity.useCell = 13
			case 10:
				ps.playerEntity.useCell = 14
			case 11:
				ps.playerEntity.useCell = 13
			case 12:
				ps.playerEntity.useCell = 12
			case 13:
				ps.playerEntity.useCell = 11
			case 14:
				ps.playerEntity.useCell = 10
			case 15:
				ps.playerEntity.useCell = 9
			case 16:
				ps.playerEntity.useCell = 8
			case 17:
				ps.playerEntity.useCell = 7
			case 18:
				ps.playerEntity.useCell = 6
			case 19:
				ps.playerEntity.useCell = 5
			case 20:
				ps.playerEntity.useCell = 4
			case 21:
				ps.playerEntity.useCell = 3
			}
		} else {
			ps.playerEntity.useCell = 0
		}

		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell + ps.playerEntity.useCell)
	}

	if engo.Input.Button("Jump").JustPressed() {
		if ps.playerEntity.ifJumping {
			if ps.playerEntity.jumpCount < MaxCount/2 {
				ps.playerEntity.jumpCount2Step = ps.playerEntity.jumpCount - 1
			} else {
				ps.playerEntity.jumpCount2Step = MaxCount - (ps.playerEntity.jumpCount - 1)
			}
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = false
		}

		if ps.playerEntity.jumpCount == 0 {
			ps.playerEntity.jumpCount2Step = 0
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = true
		}
		ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step
	}

	if ps.playerEntity.jumpCount != 0 {
		ps.playerEntity.jumpCount++
		if ps.playerEntity.jumpCount <= ps.playerEntity.topCount {
			ps.playerEntity.SpaceComponent.Position.Y -= JumpHeight
		} else if ps.playerEntity.jumpCount <= ps.playerEntity.bottomCount {
			ps.playerEntity.SpaceComponent.Position.Y += JumpHeight
		} else {
			ps.playerEntity.jumpCount = 0
			ps.playerEntity.ifJumping = false
		}
	}

	if ps.playerEntity.jumpCount == 0 {
		if getMakingInfo(FallPoint, int(ps.playerEntity.LeftPositionX)) && getMakingInfo(FallPoint, int(ps.playerEntity.RightPositionX)) {
			ps.playerEntity.ifFalling = true
			ps.playerEntity.SpaceComponent.Position.Y += MoveDistance
		}
	}
	if ps.playerEntity.SpaceComponent.Position.Y > engo.WindowHeight() {
		ps.Remove(ps.playerEntity.BasicEntity)
	}
	if ps.playerEntity.ifFalling {
		return
	}
}

func (ps *PlayerSystem) Remove(ecs.BasicEntity) {
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(ps.playerEntity.BasicEntity)
		}
	}
}
