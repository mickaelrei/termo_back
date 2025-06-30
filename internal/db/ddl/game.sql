-- DDL to create the game table
CREATE TABLE IF NOT EXISTS game (
    id         INTEGER  NOT NULL PRIMARY KEY AUTO_INCREMENT,
    id_user    INTEGER  NOT NULL,
    is_active  BOOLEAN  NOT NULL DEFAULT TRUE,
    FOREIGN KEY (id_user) REFERENCES user (id)
);

-- DDL to create the game word table
CREATE TABLE IF NOT EXISTS game_word (
    id_game INTEGER NOT NULL,
    word    TEXT    NOT NULL,
    idx     INTEGER NOT NULL,
    FOREIGN KEY (id_game) REFERENCES game (id)
);


-- DDL to create the game attempt table
CREATE TABLE IF NOT EXISTS game_attempt (
    id_game INTEGER NOT NULL,
    attempt TEXT    NOT NULL,
    idx     INTEGER NOT NULL,
    FOREIGN KEY (id_game) REFERENCES game (id)
);
