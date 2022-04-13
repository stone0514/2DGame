package systems

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	tileFile = "./tileSet/world/tileSpritesheet.png"
	//FallPoint ...
	FallPoint []int
	//PipePoint ...
	PipePoint []int
	//0: Not Making, 1:Making, 2:Other
	makingFall      int
	makingCloud     int
	makingPipe      int
	addCell         int
	pipePositionY   float32
	onPipePositionY float32
)

//Tile ...
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

//TileSystem ...
type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
}

//Update ...
//update is ran every frame, with "dt" being the time
func (*TileSystem) Update(dt float32) {}

//Remove ...
//remove is called whenever an entity is removed from the world ...
func (*TileSystem) Remove(ecs.BasicEntity) {}

//New ...
//initialisation of the system
func (ts *TileSystem) New(w *ecs.World) {
	ts.world = w
	//create tile spritesheet
	Spritesheet33x33 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth33, CellHeight33, 0, 0)
	//create cloud spritesheet
	Spritesheet98x66 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth98, CellHeight60, 0, 0)
	//create pipe spritesheet
	Spritesheet66x66 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth66, CellHeight66, 0, 0)

	addCell = 0
	cloudHeight := 0

	//initialize pipePosition
	pipePositionY = engo.WindowHeight() - CellHeight33*6
	onPipePositionY = engo.WindowHeight() - CellHeight33*8

	makingFall = 0

	Tiles := make([]*Tile, 0)

	for i := 0; i <= TileNum; i++ {
		//create ground
		if i >= 10 && i < TileNum-AroundGoalTileNum {
			randomNum := rand.Intn(10)
			if randomNum == 0 {
				makingFall = 1
			} else {
				if makingFall == 1 {
					makingFall = 2
				} else {
					makingFall = 0
				}
			}
		}
		if makingFall != 0 {
			for j := 0; j < CellWidth33; j++ {
				//fallpoint position record
				FallPoint = append(FallPoint, i*CellWidth33+j)
			}
		} else {
			for j := 0; j < TileDepth; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				//spaceComponent
				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth33), Y: float32(int(engo.WindowHeight()) - (j+1)*CellHeight33)},
				}
				//renderComponent
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet33x33.Cell(GroundSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				//set component
				Tiles = append(Tiles, tile)
			}
		}
		//create cloud
		if makingCloud == 0 {
			randomNum := rand.Intn(12)
			if randomNum < 3 {
				makingCloud = 1
				cloudHeight = randomNum
			}
		}
		if makingCloud != 0 {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			j := float32(0)
			//creating second cloud
			if makingCloud > 2 {
				j = float32(i) - 0.7
			} else {
				j = float32(i)
			}

			//spaceComponent
			tile.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(j * CellWidth98), Y: float32(int(engo.WindowHeight()/3) - cloudHeight*CellHeight60)},
			}

			//renderComponent
			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet98x66.Cell(CloudSpriteSheetCell),
				Scale:    engo.Point{X: 1, Y: 1},
			}
			tile.RenderComponent.SetZIndex(float32(makingCloud))

			//set component
			Tiles = append(Tiles, tile)

			switch makingCloud {
			case 1:
				makingCloud++
				addCell = 1
				break
			case 2:
				makingCloud++
				addCell = 1
				break
			default:
				makingCloud = 0
				addCell = 0
				break
			}
		}
	}

	//create pipePosition
	for i := 0; i <= TileNum; i++ {
		makingPipe = 1
		for j := 0; j < PipeTileNum+2; j++ {
			if GetMakingInfo(FallPoint, (i+j)*CellWidth33) {
				makingPipe = 0
			}
		}
		if i >= 10 && i < TileNum-AroundGoalTileNum {
			if makingPipe != 0 {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				//spaceComponent
				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth33), Y: pipePositionY},
				}

				//renderComponent
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet66x66.Cell(PipeSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				//set component
				Tiles = append(Tiles, tile)

				//record pipePosition
				for j := 0; j < CellWidth66; j++ {
					PipePoint = append(PipePoint, i*CellWidth33+j)
				}
				//increment
				i = i + 10
			}
		}
	}

	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}
