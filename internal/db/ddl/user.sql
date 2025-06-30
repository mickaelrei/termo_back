-- DDL to create user table
CREATE TABLE IF NOT EXISTS user (
    id       INTEGER     NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name     VARCHAR(32) NOT NULL,
    password TEXT        NOT NULL,
    score    INTEGER     NOT NULL DEFAULT 0,
    UNIQUE KEY (name)
);
