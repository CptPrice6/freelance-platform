-- +goose Up
ALTER TABLE refresh_tokens
  ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE refresh_tokens DROP CONSTRAINT fk_user_id;
