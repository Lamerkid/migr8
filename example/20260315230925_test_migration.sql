-- +migr8:up

CREATE TABLE IF NOT EXISTS test (
    id UUID PRIMARY KEY
);

-- +migr8:down

DROP TABLE IF EXISTS test;
