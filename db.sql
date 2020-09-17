CREATE DATABASE Truth_or_dare;
USE Truth_or_dare;
-- modify the telegram_id field to be unique. Run this query kwa the db

CREATE TABLE players
(
    user_id     INT(10) PRIMARY KEY AUTO_INCREMENT,
    telegram_id INT(15)   NOT NULL,
    first_name  VARCHAR(30) NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE game_session
(
    game_id    INT(15) PRIMARY KEY AUTO_INCREMENT,
    active     BOOL,
    updated_at DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);



CREATE TABLE player_scores
(
    player_score_id INT(30) PRIMARY KEY AUTO_INCREMENT,
    user_id         INT,
    game_id         INT,
    scores          INT(30) DEFAULT 0,

    FOREIGN KEY (user_id) REFERENCES players (user_id),
    FOREIGN KEY (game_id) REFERENCES game_session (game_id)
);

CREATE TABLE truths_dares (
    id        INT(10) PRIMARY KEY AUTO_INCREMENT,
    challenge VARCHAR(256) NOT NULL,
    type      ENUM ('truth', 'dare')
);

INSERT INTO truths_dares (challenge, type)
VALUES
('What are your top three turn-ons?', 'truth'),
('What is your deepest darkest fear?', 'truth'),
('Tell me about your first kiss', 'truth'),
('Who is the sexiest person here?', 'truth'),
('What is your biggest regret?', 'truth'),
('Who here has the nicest butt?', 'truth'),
('Who is your crush?', 'truth'),
('Who was the last person you licked?', 'truth');

INSERT INTO truths_dares (challenge, type)
VALUES
('Kiss the person to your left', 'dare'),
('Attempt to do a magic trick', 'dare'),
('Do four cartwheels in row', 'dare'),
('Let someone shave part of your body', 'dare'),
('Eat five tablespoons of a condiment', 'dare');




create unique index unique_telegram_id on  players (telegram_id);
alter table game_session add column number_of_players int default 1;