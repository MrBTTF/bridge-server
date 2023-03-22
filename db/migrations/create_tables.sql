DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS players CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;


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
    ADD FOREIGN KEY (current_player) REFERENCES users (user_id) ON DELETE CASCADE;
    
ALTER TABLE players
    ADD FOREIGN KEY (session_id) 
        REFERENCES sessions (session_id) ON DELETE CASCADE,
    ADD FOREIGN KEY (user_id) 
        REFERENCES users (user_id) ON DELETE CASCADE;

ALTER TABLE rooms
    ADD FOREIGN KEY (host_id) REFERENCES users (user_id) ON DELETE CASCADE; 

GRANT ALL ON ALL TABLES IN SCHEMA public TO bridge;


INSERT INTO users (user_id, email, password, nickname, token)
values ('6e3f9165-3daf-4fdc-9b52-12e6fdd810c1', 
        'test1@bridge.test', '42f749ade7f9e195bf475f37a44cafcb', 'User1', ''),  
       ('a1983803-7eb9-477a-a860-9a652adfa30d', 
       'zalizniak@zalizniak', 'ac0ddf9e65d57b6a56b2453386cd5db5', 'zalizniak', ''),
       ('be2d3a8e-5f95-4716-a0ef-a814b89dbabc', 
       'test2@bridge.test', '42f749ade7f9e195bf475f37a44cafcb', 'SecondUser', ''),
       ('851e1367-610f-412f-a840-4dfe0d9db38d', 
       'test3@bridge.test', '42f749ade7f9e195bf475f37a44cafcb', 'OkUser123', ''),
       ('af58fe77-a6bb-4169-a960-3107d7bea057', 
       'test4@bridge.test', '42f749ade7f9e195bf475f37a44cafcb', 'LastUser1', '');


