CREATE TABLE IF NOT EXISTS queue (
    version INT,
    added   INT,
    url     VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS in_progress (
    version INT,
    added   INT,
    url     VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS finished (
    version     INT,
    added       INT,
    finished    INT,
    url         VARCHAR PRIMARY KEY
);