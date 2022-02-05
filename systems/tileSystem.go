package systems

import (
	"fmt"
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

func (*TileSystem) Update(dt float32) {}

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) New(w *ecs.World) {
	ts.world = w
	SpriteSheet16x16 := common.NewSpriteSheetWithBorderFromFile(tileFile, CellWidth16, CellHeight16, 0, 0)

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
					Position: engo.Point{X: float32(i * CellWidth16), Y: float32(int(engo.WindowHeight()) - (j+i)*CellHeight16)},
				}

				tile.RenderComponent = common.RenderComponent{
					Position: engo.SpriteSheet16x16.Cell(GroundSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)
			}
		}
	}

	fmt.Println("TileSystem was added to the Scene")
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
}

const (
	TileNum               = 200
	GoalTileNum           = 10
	AroundGoalTileNum     = 30
	MountainTileNum       = 5
	PipeTileNum           = 1
	CellWidth16           = 16
	CellHeight16          = 16
	CellWidth32           = 32
	CellHeight32          = 32
	CellHeight64          = 64
	TileDepth             = 4
	GroundSpriteSheetCell = 0
	CloudSpriteSeetCell   = 0
	MountSpriteSeetCell   = 0
	PipeSpriteSheetCell   = 0
)

var tileFile = "pixel_platform_01_tileset_final.png"

var FallPoint []int
var MountainPoint []int
var PipePoint []int

//0: Not Making, 1:Making, 2:Other
var makingFall int
var makingCloud int
var makingMount int
var makingPipe int

var addCell int

var mountainPositionY float32
var pipePositionY float32
var onPipePositionY float32
