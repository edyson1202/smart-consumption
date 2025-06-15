-- +goose Up
-- +goose StatementBegin
CREATE TABLE datapoint(
                        id SERIAL PRIMARY KEY,
                        user_id INT REFERENCES users (id) ON DELETE CASCADE,
                        measured_at TIMESTAMP,
                        temperature FLOAT,
                        humidity FLOAT,
                        consumption FLOAT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE datapoint;
-- +goose StatementEnd