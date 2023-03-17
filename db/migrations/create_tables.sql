DROP TABLE sessions CASCADE;
DROP TABLE players CASCADE;
DROP TABLE users CASCADE;

CREATE TABLE IF NOT EXISTS sessions (
    session_id text PRIMARY KEY,
    players    text NOT NULL,
    deck    text[][],
    session_table    text[][],
	current_player text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS players (
    user_id text NOT NULL,
    nickname text,
    cards    text[][],
    state    smallint,
    state_name    text,
    session_id   text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, session_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id text PRIMARY KEY,
    email text,
    password text,
    nickname text,
    token text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    room_id text PRIMARY KEY,
    host_id text,
    user_ids    text[][],
    open        boolean,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE sessions
    ADD FOREIGN KEY (current_player) REFERENCES users (user_id);
    
ALTER TABLE players
    ADD FOREIGN KEY (session_id) 
        REFERENCES sessions (session_id),
    ADD FOREIGN KEY (user_id) 
        REFERENCES users (user_id);

ALTER TABLE rooms
    ADD FOREIGN KEY (host_id) REFERENCES users (user_id);
    

INSERT INTO users (user_id, email, password, nickname, token)
values ('6e3f9165-3daf-4fdc-9b52-12e6fdd810c1', 'test1@bridge.test', '482c811da5d5b4bc6d497ffa98491e38', 'Test User', ''),
       ('a1980803-7eb9-477a-a860-9a652adfa30d', 'zalizniak@zalizniak', 'ac0ddf9e65d57b6a56b2453386cd5db5', 'zalizniak', '');
