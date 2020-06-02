CREATE TABLE IF NOT EXISTS organisations
(
    id   SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS teams
(
    id              SERIAL PRIMARY KEY,
    organisation_id INT          not null REFERENCES organisations (id) ON DELETE CASCADE,
    name            VARCHAR(255) not null,
    UNIQUE (organisation_id, name)
);

CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    organisation_id INT          NOT NULL REFERENCES organisations (id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    password        VARCHAR(255) NOT NULL,
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS tasks
(
    id          CHAR(36) PRIMARY KEY,
    author_id   INT          NOT NULL REFERENCES users (id),
    team_id     INT          NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    title       VARCHAR(70)  NOT NULL,
    description TEXT,
    status      VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks_events
(
    id          SERIAL PRIMARY KEY,
    task_id     CHAR(36)     NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    event_type  VARCHAR(255) NOT NULL,
    payload     TEXT,
    occurred_on TIMESTAMP    NOT NULL
);

CREATE TABLE IF NOT EXISTS task_comments
(
    id         SERIAL PRIMARY KEY,
    task_id    CHAR(36)  NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    author_id  INT       REFERENCES users (id) ON DELETE SET NULL,
    message    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS task_last_watched_comment
(
    user_id         INT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    task_id         CHAR(36) NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    last_comment_id INT      NOT NULL REFERENCES task_comments (id),
    UNIQUE (user_id, task_id)
);

CREATE TABLE IF NOT EXISTS sub_task_stages
(
    id      SERIAL PRIMARY KEY,
    team_id INT NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    rank    INT NOT NULL,
    name    VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS sub_tasks
(
    id         SERIAL PRIMARY KEY,
    task_id    CHAR(36)     NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    stage_id   INT          NOT NULL REFERENCES sub_task_stages (id),
    author_id  INT          NOT NULL REFERENCES users (id),
    status     INT          NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    closed_at  TIMESTAMP
);



