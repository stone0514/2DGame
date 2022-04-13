package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var playerFile = "./characters/robot.png"

//Player ...
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	playerPositionY float32
	//step position
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
	ifOnPipe           bool
	ifFalling          bool
	ifStart            bool
}

//PlayerSystem ...
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
}

//New ...
func (ps *PlayerSystem) New(w *ecs.World) {
	//add world
	ps.world = w
	//generate Entity
	player := Player{BasicEntity: ecs.NewBasic()}

	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 3200, Y: 3200},
	}

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
	//renderComponent
	player.RenderComponent = common.RenderComponent{
		Drawable: player.spritesheet.Cell(PlayerSpriteSheetCell),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	player.RenderComponent.SetZIndex(5)

	//set component
	ps.playerEntity = &player

	//initialize
	ps.playerEntity.playerPositionY = PsPositionY
	ps.playerEntity.LeftPositionX = PsPositionX + float32(ExtraSizeX)
	ps.playerEntity.RightPositionX = PsPositionX + CellWidth32 - float32(ExtraSizeX)
	ps.playerEntity.ifFalling = false
	ps.playerEntity.ifOnPipe = false
	ps.playerEntity.cameraMoveDistance = 0
	ps.playerEntity.topCount = 1 + MaxCount/2
	ps.playerEntity.bottomCount = 0

	//add renderSystem
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
	//camera setting
	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 3200, Y: 300},
	}
}

//Update ...
func (ps *PlayerSystem) Update(dt float32) {

	//no move at rest
	if ps.playerEntity.SpaceComponent.Position.Y == ps.playerEntity.playerPositionY {
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell)
	}
	//no move on the pipe
	if ps.playerEntity.ifOnPipe {
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell)
	}

	//return if you die
	if ifGameOver {
		return
	}
	//return not start
	if !ps.playerEntity.ifStart {
		return
	}

	//Don't move to the right when reach the goal
	if int(ps.playerEntity.LeftPositionX) >= (TileNum-GoalTileNum+2)*CellWidth33 {
		ps.Remove(ps.playerEntity.BasicEntity)
	}

	//fall fallpoint judgment
	if ps.playerEntity.jumpCount == 0 {
		if GetMakingInfo(FallPoint, int(ps.playerEntity.LeftPositionX)) && GetMakingInfo(FallPoint, int(ps.playerEntity.RightPositionX)) {
			ps.playerEntity.ifFalling = true
			ps.playerEntity.SpaceComponent.Position.Y += MoveDistance
		}
	}
	//remove if fall and call PlayerDie()
	if ps.playerEntity.SpaceComponent.Position.Y > 500 {
		ps.Remove(ps.playerEntity.BasicEntity)
		ps.PlayerDie()
	}
	if ps.playerEntity.ifFalling {
		return
	}

	if int(ps.playerEntity.LeftPositionX) > (TileNum-GoalTileNum-3)*CellWidth33 {
		for _, system := range ps.world.Systems() {
			switch sys := system.(type) {
			case *HUDTextSystem:
				sys.TextInit(sys.TextEntity, TextGoal)
			}
		}
		ps.Remove(ps.playerEntity.BasicEntity)
	}

	//move right player
	if engo.Input.Button("MoveRight").Down() {
		//if the player is under the pipe
		if GetMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) && int(ps.playerEntity.SpaceComponent.Position.Y) > int(engo.WindowHeight())-CellHeight33*8 {
			//Can't move right
		} else {
			//on a pipe and not jumping
			if ps.playerEntity.ifOnPipe && ps.playerEntity.jumpCount == 0 {
				if !GetMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) && !GetMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
					ps.playerEntity.ifOnPipe = false
					ps.playerEntity.SpaceComponent.Position.Y = ps.playerEntity.playerPositionY
				}
			}
			//Don't move camera if it's located to the left of the center of the screen
			if int(ps.playerEntity.SpaceComponent.Position.X) < ps.playerEntity.cameraMoveDistance+int(engo.WindowWidth())/2 {
				ps.playerEntity.SpaceComponent.Position.X += MoveDistance
				ps.playerEntity.LeftPositionX += MoveDistance
				ps.playerEntity.RightPositionX += MoveDistance
			} else {
				//move player if not reaching the right edge of the screen
				if int(ps.playerEntity.SpaceComponent.Position.X) < int(engo.WindowWidth())-CellWidth32 {
					ps.playerEntity.SpaceComponent.Position.X += MoveDistance
					ps.playerEntity.LeftPositionX += MoveDistance
					ps.playerEntity.RightPositionX += MoveDistance
				}
				if int(ps.playerEntity.SpaceComponent.Position.X) < TileNum*CellWidth33-int(engo.WindowWidth())/2 {
					//move camera
					engo.Mailbox.Dispatch(common.CameraMessage{
						Axis:        common.XAxis,
						Value:       MoveDistance,
						Incremental: true,
					})
				}
				ps.playerEntity.cameraMoveDistance += MoveDistance
			}
		}
		//not jumping animation
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

		//change player behavior
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell + ps.playerEntity.useCell)
	}

	//jump
	if engo.Input.Button("Jump").JustPressed() {
		//2stepJump
		if ps.playerEntity.ifJumping {
			if ps.playerEntity.jumpCount < MaxCount/2 {
				ps.playerEntity.jumpCount2Step = ps.playerEntity.jumpCount - 1
			} else {
				ps.playerEntity.jumpCount2Step = MaxCount - (ps.playerEntity.jumpCount - 1)
			}
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = false
		}

		//first jump
		if ps.playerEntity.jumpCount == 0 {
			ps.playerEntity.jumpCount2Step = 0
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = true
		}
		//jumping on the pipe
		if ps.playerEntity.ifOnPipe {
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step + 11
			//jumping on the ground
		} else {
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step
		}
	}

	if ps.playerEntity.jumpCount != 0 {
		ps.playerEntity.jumpCount++
		if ps.playerEntity.jumpCount <= ps.playerEntity.topCount {
			ps.playerEntity.SpaceComponent.Position.Y -= JumpHeight
		} else if ps.playerEntity.jumpCount <= ps.playerEntity.bottomCount {
			if ps.playerEntity.SpaceComponent.Position.Y == onPipePositionY {
				//if right or left foot is on the pipe
				if GetMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || GetMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
					ps.playerEntity.jumpCount = 0
					ps.playerEntity.ifJumping = false
					ps.playerEntity.ifOnPipe = true
				} else {
					//down
					ps.playerEntity.SpaceComponent.Position.Y += JumpHeight
				}
			} else {
				//down
				ps.playerEntity.SpaceComponent.Position.Y += JumpHeight
			}
		} else {
			ps.playerEntity.jumpCount = 0
			ps.playerEntity.ifJumping = false
			//when the landing point is a pipe
			if GetMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || GetMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
				ps.playerEntity.ifOnPipe = true
				//when the landing point is a ground
			} else {
				ps.playerEntity.ifOnPipe = false
			}
		}
	}
}

//Remove removes an Entity from the System
func (ps *PlayerSystem) Remove(ecs.BasicEntity) {
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(ps.playerEntity.BasicEntity)
		}
	}
}
