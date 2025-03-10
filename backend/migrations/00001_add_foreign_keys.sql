-- +goose Up
ALTER TABLE refresh_tokens
  ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE client_data
  ADD CONSTRAINT fk_client_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE freelancer_data
  ADD CONSTRAINT fk_freelancer_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE freelancer_skills
  ADD CONSTRAINT fk_freelancer FOREIGN KEY (freelancer_data_id) REFERENCES freelancer_data(id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_skill FOREIGN KEY (skills_id) REFERENCES skills(id) ON DELETE CASCADE;




-- +goose Down
ALTER TABLE refresh_tokens DROP CONSTRAINT fk_user_id;
ALTER TABLE client_data DROP CONSTRAINT fk_client_user;
ALTER TABLE freelancer_data DROP CONSTRAINT fk_freelancer_user;

ALTER TABLE freelancer_skills 
DROP CONSTRAINT fk_freelancer, 
DROP CONSTRAINT fk_skill;
