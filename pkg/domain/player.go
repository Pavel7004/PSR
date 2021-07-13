package domain

type Player struct {
	id string
}

func NewPlayer(id string) *Player {
	return &Player{
		id: id,
	}
}

func (p *Player) GetID() string {
	return p.id
}
