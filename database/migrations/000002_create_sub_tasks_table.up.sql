CREATE TABLE IF NOT EXISTS sub_task_stages
(
    id   SERIAL PRIMARY KEY,
    rank INT NOT NULL,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS sub_tasks
(
    id         SERIAL PRIMARY KEY,
    task    CHAR(36)     NOT NULL REFERENCES tasks(id),
    stage      INT          NOT NULL REFERENCES sub_task_stages (id),
    author     INT          NOT NULL REFERENCES users (id),
    status     INT          NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    closed_at  TIMESTAMP
);