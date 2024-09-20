CREATE TABLE signs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_id VARCHAR(50) NOT NULL,
    discord_id VARCHAR(50) NOT NULL
);

CREATE TABLE messages (
    message_id     VARCHAR(100) NOT NULL PRIMARY KEY,
    week           VARCHAR(2) NOT NULL,
    year           VARCHAR(4) NOT NULL,
    game_1         VARCHAR(50) NOT NULL,
    game_2         VARCHAR(50) NOT NULL,
    game_3         VARCHAR(50) NOT NULL,
    creation_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shuffled_teams (
    id INT AUTO_INCREMENT PRIMARY KEY,
    shuffle_id VARCHAR(50),
    team TINYINT NOT NULL,
    discord_id VARCHAR(50) NOT NULL,
    creation_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);