CREATE TABLE IF NOT EXISTS shortkeys (
    id integer NOT NULL,
    shortkey text,
    taken boolean,
    PRIMARY KEY(id)
);