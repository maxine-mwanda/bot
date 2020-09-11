package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"telegrambot/resources/db/utils"
)

// TODO: When someone chooses Three, create a game session

type game_session struct {
	game_id int `json:"game_id"`
	active int `json:"active"`
}

type player_scores struct {
	user_id int `json:"user_id"`
	game_id int `json:"game_id"`
}

func CreateGameSession() (id int64, err error) {
	query := "insert into game_session (active) values ('1')"
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect to db")
		return
	}

	row, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, _ = row.LastInsertId()
	return
}

func FetchGameId() (game_id int64, err error) {
	query := "insert into game_session (active) values (?)"
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect todb")
		return
	}
	row, err := db.Exec(query, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	game_id, _ = row.LastInsertId()
	return
}

	//var gs game_session
	//for rows.Next() {
	//	if err = rows.Scan(&gs.game_id, &gs.active, &game_session.updated_at); err != nil {
	//		fmt.Println("Unable to scan because ", err)
	//		continue
	//	}
	//	game_session = append(game_session, err)
	//}
	//return
//}

func JsonResponse(writer http.ResponseWriter, status int, response interface{}) {
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonResponse)
}


func Scores () (user_id int64, err error){
query := "insert into player_scores (user_id) values ('567')"
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect to db")
		return
	}

	row, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	user_id, _ = row.LastInsertId()
	return Scores()
}


func Create_player(telegramId int, firstName string) (id int64, err error) {
	query := "insert into players (telegram_id, first_name)" +
		"values (?,?)"
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect todb")
		return
	}
	row, err := db.Exec(query, telegramId, firstName)
	if err != nil {
		// TODO: this means the player exists - "select user_id from players where telegram_id=?   ... then return the user_id
		query = "select user_id from players where telegram_id=?"
		playerRow := db.QueryRow(query, telegramId)
		if playerRow == nil {
			log.Println("An error occured. Player does not exist and could not be created")
			return
		}
		err = playerRow.Scan(&id)
		if err != nil {
			log.Println("Unable to scan telegram_id because", err)
			return
		}
		return
	}

	id, _ = row.LastInsertId()
	return
}

func PlayerScores(user_id, game_id int)  (err error){
	query := "insert into player_scores (user_id, game_id, scores)" +
		"values (?,?,?)"
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect todb")
		return
	}
	_, err = db.Exec(query, user_id, game_id, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}