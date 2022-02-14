package systems

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	TileNum           = 200
	GoalTileNum       = 10
	AroundGoalTileNum = 30
	PipeTileNum       = 2
	// tileSize 33x33
	CellWidth33  = 33
	CellHeight33 = 33
	// cloudSize 99x60
	CellWidth99  = 99
	CellHeight60 = 60
	// PipeSize
	CellWidth66  = 66
	CellHeight66 = 66

	CellWidth32  = 32
	CellHeight32 = 32
	CellHeight64 = 64
	TileDepth    = 4
	// tileSpriteSheet 28x27
	GroundSpriteSheetCell = 133
	// cloudSpriteSheet 10x16?
	CloudSpriteSheetCell = 140
	// pipeSpriteSheet
	PipeSpriteSheetCell = 120
)

var (
	tileFile = "./tileSet/world/tileSpritesheet.png"

	FallPoint []int
	PipePoint []int

	//0: Not Making, 1:Making, 2:Other
	makingFall  int
	makingCloud int
	makingPipe  int

	addCell int

	pipePositionY   float32
	onPipePositionY float32
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
}

func (*TileSystem) Update(dt float32) {}

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) New(w *ecs.World) {
	ts.world = w
	Spritesheet33x33 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth33, CellHeight33, 0, 0)
	Spritesheet99x66 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth99, CellHeight60, 0, 0)
	Spritesheet66x66 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth66, CellHeight66, 0, 0)

	addCell = 0
	cloudHeight := 0

	pipePositionY = engo.WindowHeight() - CellHeight33*6
	onPipePositionY = engo.WindowHeight() - CellHeight33*8

	makingFall = 0

	Tiles := make([]*Tile, 0)

	for i := 0; i <= TileNum; i++ {
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
				FallPoint = append(FallPoint, i*CellWidth33+j)
			}
		} else {
			for j := 0; j < TileDepth; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth33), Y: float32(int(engo.WindowHeight()) - (j+1)*CellHeight33)},
				}
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet33x33.Cell(GroundSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)
			}
		}
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
			if makingCloud > 2 {
				j = float32(i) - 0.7
			} else {
				j = float32(i)
			}

			tile.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(j * CellWidth99), Y: float32(int(engo.WindowHeight()/3) - cloudHeight*CellHeight60)},
			}

			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet99x66.Cell(CloudSpriteSheetCell),
				Scale:    engo.Point{X: 1, Y: 1},
			}
			tile.RenderComponent.SetZIndex(float32(makingCloud))

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

	for i := 0; i <= TileNum; i++ {
		makingPipe = 1
		for j := 0; j < PipeTileNum+2; j++ {
			if getMakingInfo(FallPoint, (i+j)*CellWidth33) {
				makingPipe = 0
			}
		}
		if i >= 10 && i < TileNum-AroundGoalTileNum {
			if makingPipe != 0 {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth33), Y: pipePositionY},
				}

				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet66x66.Cell(PipeSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)

				for j := 0; j < CellWidth66; j++ {
					PipePoint = append(PipePoint, i*CellWidth33+j)
				}
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

func getMakingInfo(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
