package resources

import (
	"encoding/json"
	"log"
	"net/http"
	"telegrambot/resources/db/utils"
)


func CreateGameSession() (id int64, err error) {
	query := "insert into game_session (Active) values ('1')"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect to db")
		return
	}

	row, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	id, _ = row.LastInsertId()
	return
}

func FetchGameId() (game_id int64, err error) {
	query := "insert into game_session (active) values (?)"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect todb")
		return
	}
	row, err := db.Exec(query, 1)
	if err != nil {
		log.Println(err)
		return
	}
	game_id, _ = row.LastInsertId()
	return

}

func JsonResponse(writer http.ResponseWriter, status int, response interface{}) {
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonResponse)
}


func Scores (gameId string, userId int64) (err error){
query := "insert into player_scores (user_id, game_id) values (?, ?)"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect to db")
		return
	}

	_, err = db.Exec(query, userId, gameId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("created record in player_scores")
	return
}


func Create_player(telegramId int, firstName string) (id int64, err error) {
	query := "insert into players (telegram_id, first_name)" +
		"values (?,?)"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect todb")
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

func PlayerScores(telegramId int, gameId string) (err error) {
	query := "update player_scores set scores = scores + 10 where game_id=? " +
		"and user_id=(select user_id from players where telegram_id=?);"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect to db")
		return
	}
	_, err = db.Exec(query, gameId, telegramId)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func GetFinalScore(gameId string, telegramId int) (score int64, err error){
	query := "select scores from player_scores where game_id=? and user_id=(select user_id from players where telegram_id=?)"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect to db")
		return
	}

	if err = db.QueryRow(query, gameId, telegramId).Scan(&score); err != nil {
		log.Println("Unable to fetch score because ", err)
		return
	}
	log.Println("Found score as ", score)
	return
}

func SetGamePlayers(gameId, numberOfPlayers string) (err error) {
	query := "update game_session set number_of_players=? where game_id=?"
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect todb")
		return
	}
	if _, err = db.Exec(query, numberOfPlayers, gameId); err != nil {
		log.Println("Unable to update the game session because ", err)
		return
	}
	log.Println("Game session updated successfully")
	return
}

func SpaceAvailableInGameSession(gameId string) (ok bool, err error) {
	db, err := utils.Connecttodb()
	if err != nil {
		log.Println("unable to connect todb")
		return
	}

	queryMaxPlayers := "select number_of_players from game_session where game_id=? and active=1"
	var maxNumberOfPlayers int
	if err = db.QueryRow(queryMaxPlayers, gameId).Scan(&maxNumberOfPlayers); err != nil {
		log.Println("Unable to read maximum players because ", err)
		return
	}

	var joinedNumberOfPlayers int
	queryNoJoinedPlayers := "select count(*) from player_scores where game_id=?"
	if err = db.QueryRow(queryNoJoinedPlayers, gameId).Scan(&joinedNumberOfPlayers); err != nil {
		log.Println("Unable to count joined players because ", err)
		return
	}
	ok = joinedNumberOfPlayers < maxNumberOfPlayers
	log.Println("Max : ", maxNumberOfPlayers, " Joined : ", joinedNumberOfPlayers)
	return
}