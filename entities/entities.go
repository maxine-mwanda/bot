package entities

type TruthAndDare struct {
	ID int64 `json:"id"`
	Challenge string `json:"challenge"`
	Type string `json:"type"`
}

type GameSession struct {
	GameId int `json:"game_id"`
	Active int `json:"active"`
}

type PlayerScores struct {
	UserId int `json:"user_id"`
	GameId int `json:"game_id"`
}