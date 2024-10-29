DROP TABLE IF EXISTS webhooks;
CREATE TABLE IF NOT EXISTS webhooks
(
    name             VARCHAR(255) NOT NULL PRIMARY KEY,
    path             VARCHAR(255) NOT NULL,
    method           VARCHAR(10)  NOT NULL,
    description      VARCHAR(255),
    jq_filter        TEXT,
    forward_to       TEXT,
    preserve_payload BOOLEAN
);