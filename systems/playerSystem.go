package systems

import (
	"github.com/EngoEngine/ecs"
)

const (
	MoveDistance          = 4
	JumpHeight            = 4
	MaxCount              = 40
	PlayerSpriteSheetCell = 5
	ExtraSizeX            = 4
)

var playerFile = "./assets/character/Adventure/adventure-Sheet.png"

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	playerPositionY    float32
	LeftPositionX      float32
	RightPositionX     float32
	cameraMoveDistance int
	spriteSheet        *common.SpriteSheet
	useCell            int
	jumpCount          int
	jumpCount2Step     int
	topCount           int
	bottomCount        int
	ifJumping          bool
	ifOnPipe           bool
	ifFalling          bool
}

type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
}

func (ps *PlayerSystem) New(w *ecs.World) {
	ps.world = w
	player := Player{BasicEntity, ecs.NewBasic()}

	PsPositionX := float32(0)
	PsPositionY := engo.WindowHeigt() - CellHeight16*6

	playerSpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: PsPositionX, Y: PsPositionY},
		Width:    30,
		Height:   30,
	}

	player.SpriteSheet = common.NewSpriteSheetWithBorderFromFile(playerFile, 32, 32, 0, 0)
	player.RenderComponent = common.RenderComponent{
		Drawable: player.spriteSheetCell(PlayerSpriteSheetCell),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	Player.RenderComponent.SetZIndex(5)

	ps.playerEntity = &player

	ps.playerEntity.playerPositionY = PsPositionY
	ps.playerEntity.LeftPositionX = PsPositionX + float32(ExtraSizeX)
	ps.playerEntity.RightPositionX = PsPositionX + CellWidth32 - float32(ExtraSizeX)
	ps.playerEntity.ifFalling = false
	ps.playerEntity.ifOnPipe = false
	ps.playerEntity.cameraMoveDistance = 0
	ps.playerEntity.topCount = 1 + MaxCount/2
	ps.playerEntity.bottomCount = 0

	for _, system := range ps.world.System() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}

	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 3200, Y: 300},
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
			ps.playerEntity.ifJumpint = true
		}

		if ps.playerEntity.ifOnPipe {
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step + 8
		} else {
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step
		}
	}

	if ps.playerEntity.jumpCount != 0 {
		ps.playerEntity.jumpCount++
		if ps.playerEntity.jumpCount <= ps.playerEntity.topCount {
			ps.PlayerEntity.SpaceComponent.Position.Y -= jumpHeight
		} else if ps.playerEntity.jumpCount <= ps.playerEntity.bottomCount {
			if ps.playerEntity.SpaceComponent.Position.Y == onPipePositionY {
				if getMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
					ps.playerEntity.jumpCount = 0
					ps.playerEntity.ifJumping = false
					ps.playerEntity.ifOnPipe = true
				} else {
					ps.playerEntity.SpaceComponent.Position.Y += jumpHeight
				}
			} else {
				ps.playerEntity.SpaceComponent.Position.Y += jumpHeight
			}
		} else {
			ps.playerEntity.jumpCount = 0
			ps.playerEntity.ifJumping = false

			if getMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
				ps.playerEntity.ifOnPipe = true
			} else {
				ps.playerEntity.ifOnPipe = false
			}
		}
	}
}
