package domain

type Player struct {
	id    string
	score int
}

func NewPlayer(id string) *Player {
	return &Player{
		id:    id,
		score: 0,
	}
}

func (p *Player) GetID() string {
	return p.id
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) IncrementScore() {
	p.score++
}

func (p *Player) ResetScore() {
	p.score = 0
}
