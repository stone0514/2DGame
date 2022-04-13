package systems

var ifGameOver bool

//playerSet
const (
	MoveDistance          = 6
	JumpHeight            = 6
	MaxCount              = 40
	PlayerSpriteSheetCell = 0
	ExtraSizeX            = 0
)

//tileSet
const (
	//num tiles
	TileNum = 100
	//num of tiles from the final point to the goal
	GoalTileNum = 5
	//num of tiles near the goal
	AroundGoalTileNum = 15
	PipeTileNum       = 2
	//tileSize 33x33
	CellWidth33  = 33
	CellHeight33 = 33
	//cloudSize 98x60
	CellWidth98  = 98
	CellHeight60 = 60
	//pipeSize
	CellWidth66  = 66
	CellHeight66 = 66
	CellWidth32  = 32
	CellHeight32 = 32
	CellHeight64 = 64
	TileDepth    = 4
	//tileSpriteSheet 28x27
	GroundSpriteSheetCell = 133
	//cloudSpriteSheet 10x16
	CloudSpriteSheetCell = 140
	//pipeSpriteSheet
	PipeSpriteSheetCell = 120
)

//enermyType
const (
	//EnermyType0 : robot
	EnermyType0     = 0
	Type0Count      = 128
	ExtraSizeXType0 = 6
	ExtraSizeYType0 = 8
)

//HUDText
const (
	TextNone = iota
	TextTitle
	TextGoal
	TextEnd
)
