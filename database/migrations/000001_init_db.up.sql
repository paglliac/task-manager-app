CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    email    VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS tasks
(
    id          CHAR(36) PRIMARY KEY,
    author      INT          NOT NULL REFERENCES users (id),
    title       VARCHAR(70)  NOT NULL,
    description TEXT,
    status      VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks_events
(
    id          SERIAL PRIMARY KEY,
    task_id     CHAR(36)     NOT NULL REFERENCES tasks (id),
    event_type  VARCHAR(255) NOT NULL,
    payload     TEXT,
    occurred_on TIMESTAMP    NOT NULL
);

CREATE TABLE IF NOT EXISTS task_last_watched_event
(
    user_id       INT      NOT NULL REFERENCES users (id),
    task_id       CHAR(36) NOT NULL REFERENCES tasks (id),
    last_event_id INT      NOT NULL REFERENCES tasks_events (id),
    UNIQUE (user_id, task_id)
);

CREATE TABLE IF NOT EXISTS task_comments
(
    id         CHAR(36) PRIMARY KEY,
    task_id    CHAR(36)  NOT NULL REFERENCES tasks (id),
    author     INT       NOT NULL REFERENCES users (id),
    message    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
)
