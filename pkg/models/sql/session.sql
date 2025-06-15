CREATE TABLE session(
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES user (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL
)
    
ALTER TABLE session
ADD CONSTRAINT session_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

select * from session join users
                           on session.user_id = users.id;

CREATE INDEX session_token_hash_idx ON session(token_hash);

INSERT INTO session (user_id, token_hash)
VALUES(1, '12345') ON CONFLICT (user_id) DO
UPDATE
SET token_hash = '12435';