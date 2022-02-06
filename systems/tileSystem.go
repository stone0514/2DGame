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
	SpriteSheet32x32 := common.NewSpriteSheetWithBorderFromFile(tileFile, CellWidth32, CellHeight32, 0, 0)

	addCell = 0
	cloudHeight := 0

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
				j = float32(i) - 0.5
			} else {
				j = float32(i)
			}
		}

		tile.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: float32(j * CellWidth32), Y: float32(int(engo.WindowHeight()/3) - cloudHeight*CellHeight16)},
		}

		tile.RenderComponent = common.RenderComponent{
			Drawable: engo.SpriteSheet32x32.Cell(CloudSpriteSheetCell + addCell),
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

	SpriteSheet16x64 := common.NewSpriteSheetWithBorderFromFile(tileFile, CellWidth16, CellHeight64, 0, 0)

	mountPositionY = engo.WindowHeight() - CellHeight16*7

	for i := 0; i <= TileNum; i++ {
		makingMount = 1
		for j := 0; j <= MountTileNum+2; j++ {
			if getMakingInfo(FallPoint, (i+j)*CellWidth16) {
				makingMount = 0
			}
		}
		if makingMount != 0 && i < TileNum-AroundGoalTileNum {
			for j := 0; j < MountTileNum; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32((i + j) * CellWidth16), Y: mountPositionY},
				}

				tile.RenderComponent = common.RenderComponent{
					Drawable: SpriteSheet16x64.Cell(MountSpriteSheetCell + j),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)

				MountPoint = append(MountPoint, (i+j)*CellWidth16)
			}

			i = i + 20
		}
	}

	pipePositionY = engo.WindowHeight() - CellHeight16*6
	onPipePositionY = engo.WindowHeigh() - CellHeight16*8

	for i := 0; i <= TileNum; i++ {
		makingPipe = 1
		for j := 0; j <= PipeTileNum+2; j++ {
			if getMakingInfo(FallPoint, (i+j)*CellWidth16) || getMakingInfo(MountPoint, (i+j)*CellWidth16) {
				makingPipe = 0
			}
		}

		if i >= 10 && i < TileNum-AroundGoalTileNum {
			if makingPipe != 0 {
				tile = &Tile{BasicEntity: ecs.NewBasic()}

				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth16), Y: pipePositionY},
				}

				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet32x32.Cell(PipeSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}

				tile.RenderComponent.SetZIndex(0)

				Tiles = append(Tiles, tile)

				for j := 0; j <= CellWidth32; j++ {
					PipePoint = append(PipePoint, i*CellWidth16+j)
				}

				i = i + 30
			}
		}
	}

	caslePositionY = engo.WindowHeight() - CellHeight16*9

	tile := &Tile{BasicEntity: ecs.NewBasic()}

	tile.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32((TileNum - GoalTileNum) * CellWidth16), Y: castlePositionY},
	}

	texture, err := common.LoadSprite(castleFile)
	if err != nil {
		fmt.Println("Unable to Load Texture: " + castleFile + ":" + err.Error())
	}
	tile.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	tile.RenderComponent.SetZIndex(0)

	Tiles = append(Tiles, tile)

	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v, BasicEntity, &v.RenderComponent, &v.SpaceComponent)
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
	MountTileNum          = 5
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
var MountPoint []int
var PipePoint []int

//0: Not Making, 1:Making, 2:Other
var makingFall int
var makingCloud int
var makingMount int
var makingPipe int

var addCell int

var mountPositionY float32
var pipePositionY float32
var onPipePositionY float32
