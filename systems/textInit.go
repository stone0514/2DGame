package systems

import (
	"image/color"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

//TextInit ...
func (h *HUDTextSystem) TextInit(text *Text, textNo int) {
	//initialize
	text.textNo = textNo
	text.ifMaking = true
	//spaceComponent
	TextPositionX := (float32)(0)
	//TextHeight
	TextPositionY := engo.WindowHeight() - 300
	//TextSize
	size := float64(50)
	//renderComponent
	textDisplay := ""

	switch textNo {
	case TextTitle:
		//16space
		textDisplay = "                GAME START! "
		//22pace
	case TextGoal:
		textDisplay = "                      GOAL! "
	case TextEnd:
		textDisplay = "                GAME OVER! "
	case TextNone:
		textDisplay = ""
	}

	//spaceComponent
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: TextPositionX, Y: TextPositionY},
	}

	//renderComponent
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.White,
		Size: size,
	}
	fnt.CreatePreloaded()

	text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: textDisplay,
	}

	text.SetShader(common.TextHUDShader)
	text.RenderComponent.SetZIndex(10)

	//add entity
	h.TextEntity = text
	//add renderSystem
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
		}
	}
}
