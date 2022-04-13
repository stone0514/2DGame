package systems

//PlayerDie ...
func (ps *PlayerSystem) PlayerDie() {
	ifGameOver = true
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *HUDTextSystem:
			sys.TextInit(sys.TextEntity, TextEnd)
		}
	}
	ps.Remove(ps.playerEntity.BasicEntity)
}
