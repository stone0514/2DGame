package systems

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

//Text ...
type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	textNo   int
	ifMaking bool
}

//HUDTextSystem ...
type HUDTextSystem struct {
	world      *ecs.World
	TextEntity *Text
}

//Remove ...
func (h *HUDTextSystem) Remove(ecs.BasicEntity) {
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(h.TextEntity.BasicEntity)
		}
	}
}

//New ...
func (h *HUDTextSystem) New(w *ecs.World) {
	h.world = w
	text := &Text{BasicEntity: ecs.NewBasic()}
	h.TextInit(text, TextTitle)
}

//Update ...
func (h *HUDTextSystem) Update(dt float32) {
	if engo.Input.Button("Enter").Down() {
		switch h.TextEntity.textNo {
		case TextTitle:
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *PlayerSystem:
					sys.playerEntity.ifStart = true
				}
			}
			//remove textTitle
			h.Remove(h.TextEntity.BasicEntity)
			h.TextEntity.textNo = TextNone
			ifGameOver = false

		case TextGoal, TextEnd:
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *PlayerSystem:
					//retry
					sys.PlayerInit(sys.playerEntity)
					//time.Sleep 60/ms
					time.Sleep(time.Millisecond * 60)
				}
			}
			//show title
			h.TextInit(h.TextEntity, TextTitle)
		}
	}
}
