package entities

type TruthAndDare struct {
	ID int64 `json:"id"`
	Challenge string `json:"challenge"`
	Type string `json:"type"`
}

//type game_session struct {
//	game_id int `json:"game_id"`
//	active int `json:"active"`
//}
//
//type player_scores struct {
//	user_id int `json:"user_id"`
//	game_id int `json:"game_id"`
//}