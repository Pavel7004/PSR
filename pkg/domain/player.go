package domain

type Player struct {
	ID string
}

func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
	}
}

func (p *Player) GetID() string {
	return p.ID
}
