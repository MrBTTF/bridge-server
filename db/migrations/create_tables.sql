DROP TABLE sessions CASCADE;
DROP TABLE players CASCADE;
DROP TABLE users CASCADE;

CREATE TABLE IF NOT EXISTS sessions (
    session_id text PRIMARY KEY,
    players    text[][],
    deck    text[][],
    session_table    text[][],
	current_player text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS players (
    user_id text NOT NULL,
    name text,
    cards    text[][],
    state    smallint,
    state_name    text,
    session_id   text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, session_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id text PRIMARY KEY,
    name text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE sessions
    ADD FOREIGN KEY (current_player) REFERENCES users (user_id);
    
ALTER TABLE players
    ADD FOREIGN KEY (session_id) 
        REFERENCES sessions (session_id),
    ADD FOREIGN KEY (user_id) 
        REFERENCES users (user_id);

INSERT INTO users (user_id)
values ('user1'),
       ('user2');
