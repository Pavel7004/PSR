package libgame

type game struct {
	roomCount uint64
	rooms     []room
}

type room struct {
	playersCount uint64
	players      []playerStat
}

type playerStat struct {
	name        string
	wins        uint64
	losses      uint64
	gamesPlayed uint64
}

func init() {
	// instance := game{0, make([]room, 0)}
}
