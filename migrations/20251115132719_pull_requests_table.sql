-- +goose Up
CREATE TABLE pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(50) NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    status VARCHAR(10) NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMPTZ
);

CREATE INDEX idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX idx_pr_status ON pull_requests(status);

-- +goose Down
DROP TABLE IF EXISTS pull_requests;
