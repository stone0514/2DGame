package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	enermyFile = "./characters/robot.png"
	//touche the player
	ifTouched bool
	//EnermyPositionType0 ...
	EnermyPositionType0 []int
)

//Enermy ...
type Enermy struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	count      int
	enermyType int
}

//EnermySystem ...
type EnermySystem struct {
	world        *ecs.World
	enermyEntity []*Enermy
}

//New ...
func (es *EnermySystem) New(w *ecs.World) {
	//add world
	es.world = w
	ifTouched = false
	ifGameOver = false
	//create EnermyArray
	Enemies := make([]*Enermy, 0)

	//create spritesheet
	Spritesheet32x32 := common.NewSpritesheetWithBorderFromFile(enermyFile, 64, 66, 0, 0)

	for i := 0; i <= TileNum; i++ {
		if GetMakingInfo(PipePoint, i*CellWidth33) {
			enermy := &Enermy{BasicEntity: ecs.NewBasic()}

			//spaceComponent
			enermy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(i * CellWidth33), Y: pipePositionY},
			}

			//renderComponent
			enermy.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet32x32.Cell(7),
				Scale:    engo.Point{X: 1, Y: 1},
			}
			enermy.RenderComponent.SetZIndex(6)

			//initialize
			enermy.count = 0

			//set component
			Enemies = append(Enemies, enermy)

			//record enermy position
			for j := 0; j <= CellWidth32; j++ {
				if j > ExtraSizeXType0 && j < CellWidth32-ExtraSizeXType0 {
					EnermyPositionType0 = append(EnermyPositionType0, i*CellWidth33+j)
				}
			}
			i++
		}
	}
	//add renderSystem
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Enemies {
				es.enermyEntity = append(es.enermyEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

//Update ...
func (es *EnermySystem) Update(dt float32) {
	var playerLeftPositionX float32
	var playerRightPositionX float32
	var playerBottomPositionY float32

	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *PlayerSystem:
			if ifTouched {
				if ifGameOver {
					return
				}
				ifTouched = false
				sys.PlayerDie()
				ifGameOver = true
			}
			playerLeftPositionX = sys.playerEntity.LeftPositionX
			playerRightPositionX = sys.playerEntity.RightPositionX
			playerBottomPositionY = sys.playerEntity.SpaceComponent.Position.Y + float32(CellHeight32)
		}
	}
	if ifGameOver {
		return
	}

	for _, entity := range es.enermyEntity {
		if GetMakingInfo(EnermyPositionType0, int(playerLeftPositionX)) || GetMakingInfo(EnermyPositionType0, int(playerRightPositionX)) {
			if pipePositionY >= playerBottomPositionY && entity.SpaceComponent.Position.Y+ExtraSizeYType0 < playerBottomPositionY {
				ifTouched = true
			}
		}
		//enermyType0 animation
		if entity.enermyType == EnermyType0 {
			//Up
			if entity.count < Type0Count {
				entity.SpaceComponent.Position.Y = pipePositionY - float32(entity.count/2)
			} else if entity.count < Type0Count*2 {
				//pause
				//Down
			} else if entity.count < Type0Count*3 {
				entity.SpaceComponent.Position.Y = pipePositionY - CellHeight66 + float32((entity.count-Type0Count*2)/2)
			} else {
				entity.count = 0
			}
			entity.count++
		}
	}
}

//Remove ...
func (es *EnermySystem) Remove(ecs.BasicEntity) {}
