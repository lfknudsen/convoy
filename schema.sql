CREATE TABLE IF NOT EXISTS queue (
    version INT,
    added   INT,
    message VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS in_progress (
    version INT,
    added   INT,
    started INT,
    message VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS finished (
    version     INT,
    added       INT,
    started     INT,
    finished    INT,
    message     VARCHAR PRIMARY KEY
);