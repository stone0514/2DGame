package systems

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	TileNum               = 200
	GoalTileNum           = 10
	AroundGoalTileNum     = 30
	CellWidth16           = 33
	CellHeight16          = 33
	CellWidth32           = 32
	CellHeight32          = 32
	CellHeight64          = 64
	TileDepth             = 4
	GroundSpriteSheetCell = 133 //28x27x
)

var (
	tileFile = "./tileSet/world/tileSpritesheet.png"

	FallPoint  []int
	MountPoint []int
	PipePoint  []int

	//0: Not Making, 1:Making, 2:Other
	makingFall int

	addCell int
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
	Spritesheet16x16 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth16, CellHeight16, 0, 0)

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
			for j := 0; j < CellWidth16; j++ {
				FallPoint = append(FallPoint, i*CellWidth16+j)
			}
		} else {
			for j := 0; j < TileDepth; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth16), Y: float32(int(engo.WindowHeight()) - (j+1)*CellHeight16)},
				}
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet16x16.Cell(GroundSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)
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
