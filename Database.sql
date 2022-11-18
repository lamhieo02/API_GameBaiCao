CREATE DATABASE GameBaiCao;

CREATE TABLE users (
	username VARCHAR(30),
	coins int DEFAULT 5000, 
	point_of_3cards int DEFAULT 0,
	PRIMARY KEY (username)
)
	

CREATE TABLE decks (
	deck_id int AUTO_INCREMENT,
	remaining int,
	PRIMARY KEY(deck_id) 
)
	
CREATE TABLE cards (
	card_value int AUTO_INCREMENT,
	card_image VARCHAR(100),
	status BOOLEAN DEFAULT FALSE,
	deck_id int,
	PRIMARY KEY (card_value, deck_id),
	FOREIGN KEY(deck_id) REFERENCES decks(deck_id)
)
