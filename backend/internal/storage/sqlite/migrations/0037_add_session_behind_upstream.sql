-- +goose Up
-- +goose StatementBegin
-- behind_upstream flags when a worker session's workspace branch is out of date
-- with the upstream repository branch and Iris (The Fleet Syncer) has detected a
-- merge conflict or a pending rebase requirement.
ALTER TABLE sessions
    ADD COLUMN behind_upstream BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP COLUMN behind_upstream;
-- +goose StatementEnd
