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

type Message struct {
	CallbackQuery struct {
		From struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Data    string `json:"data"`
		Message struct {
			Chat struct {
				Id int `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	} `json:"callback_query"`

	Message struct {
		Text string `json:"text"`
		From struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Chat struct {
			Id int `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}
