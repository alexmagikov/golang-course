CREATE TABLE IF NOT EXISTS subscriptions
(
    id
    BIGSERIAL
    PRIMARY
    KEY,
    owner
    VARCHAR
(
    256
) NOT NULL,
    repo VARCHAR
(
    256
) NOT NULL,
    UNIQUE
(
    owner,
    repo
)
    );