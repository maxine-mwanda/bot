package models

import (
	"database/sql"
	"log"
	"telegrambot/entities"
)

func CreateGameSession(db *sql.DB) (id int64, err error) {
	log.Println("Creating new game session")
	query := "insert into game_session (Active) values ('1')"
	row, err := db.Exec(query)
	if err != nil {
		log.Println("unable to create game session because", err)
		return
	}
	id, _ = row.LastInsertId()
	log.Println("Created game session : ", id)
	return
}

func Scores(gameId string, userId int64, db *sql.DB) (err error) {
	query := "insert into player_scores (user_id, game_id) values (?, ?)"
	log.Println("Creating new player score")
	_, err = db.Exec(query, userId, gameId)
	if err != nil {
		log.Println("unable to create player score because", err)
		return
	}
	log.Println("created record in player_scores")
	return
}

func CreatePlayer(telegramId int, firstName string, db *sql.DB) (id int64, err error) {
	log.Println("Creating new player: ", firstName, " telegram id: ", telegramId)
	query := "insert into players (telegram_id, first_name)" +
		"values (?,?)"
	row, err := db.Exec(query, telegramId, firstName)
	if err != nil {
		log.Println("Unable to create player ", err)
		query = "select user_id from players where telegram_id=?"
		err = db.QueryRow(query, telegramId).Scan(&id)
		if err != nil {
			log.Println("Unable to scan telegram_id because", err)
			return
		}
		return
	}

	id, _ = row.LastInsertId()
	return
}

func UpdatePlayerScore(telegramId int, gameId string, db *sql.DB) (err error) {
	log.Printf("Updating score for player : %d game id : %s", telegramId, gameId)
	query := "update player_scores set scores = scores + 10 where game_id=? " +
		"and user_id=(select user_id from players where telegram_id=?);"
	_, err = db.Exec(query, gameId, telegramId)
	if err != nil {
		log.Println("Unable to update player score because", err)
		return
	}
	return
}

func GetFinalScore(gameId string, telegramId int, db *sql.DB) (score int64, err error) {
	log.Println("Fetching final score for ", telegramId)
	query := "select scores from player_scores where game_id=? and user_id=(select user_id from players where telegram_id=?)"

	if err = db.QueryRow(query, gameId, telegramId).Scan(&score); err != nil {
		log.Println("Unable to fetch final score because ", err)
		return
	}
	log.Println("Found score as ", score)
	return
}

func SetNumberOfGamePlayers(gameId, numberOfPlayers string, db *sql.DB) (err error) {
	log.Println("Setting number of players")
	query := "update game_session set number_of_players=? where game_id=?"
	if _, err = db.Exec(query, numberOfPlayers, gameId); err != nil {
		log.Println("Unable to set number of players because", err)
		return
	}
	log.Println("Number of players set successfully")
	return
}

func SpaceAvailableInGameSession(gameId string, db *sql.DB) (ok bool, err error) {
	log.Println("Checking available space in game session")
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

func TruthsAndDaresFromDB(db *sql.DB) (truthsAndDares []entities.TruthAndDare, err error) {
	log.Println("Fetching truths and dares from db")
	query := "select * from truths_dares;"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("unable to fetch truths or dares from db because", err)
		return
	}

	var t entities.TruthAndDare
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Challenge, &t.Type)
		if err != nil {
			log.Println("unable to scan from db", err)
			continue
		}
		truthsAndDares = append(truthsAndDares, t)
	}
	return
}
