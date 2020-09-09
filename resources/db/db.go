package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegrambot/resources/db/utils"
)

// TODO: When someone chooses Three, create a game session

type game_session struct {
	game_id int `json:"game_id"`
	active int `json:"active"`
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

func FetchGameId() (game_session, err error) {
	db, err := utils.Connecttodb()
	if err != nil {
		fmt.Println("unable to connect todb")
		return
	}
	query := "select * from game_session"
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("unable to get gameid")
		return
	}
	return
	//var gs game_session
	//for rows.Next() {
	//	if err = rows.Scan(&gs.game_id, &gs.active, &game_session.updated_at); err != nil {
	//		fmt.Println("Unable to scan because ", err)
	//		continue
	//	}
	//	game_session = append(game_session, err)
	//}
	//return
}

func JsonResponse(writer http.ResponseWriter, status int, response interface{}) {
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonResponse)
}