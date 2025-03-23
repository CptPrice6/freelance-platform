-- +goose Up
ALTER TABLE client_data
  ADD CONSTRAINT fk_client_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE freelancer_data
  ADD CONSTRAINT fk_freelancer_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE freelancer_skills
  ADD CONSTRAINT fk_freelancer FOREIGN KEY (freelancer_data_id) REFERENCES freelancer_data(id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_skill FOREIGN KEY (skills_id) REFERENCES skills(id) ON DELETE CASCADE;

ALTER TABLE jobs
  ADD CONSTRAINT fk_job_client FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_job_freelancer FOREIGN KEY (freelancer_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE applications
  ADD CONSTRAINT fk_application_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_application_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE;

ALTER TABLE attachments
  ADD CONSTRAINT fk_attachment_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;




-- +goose Down
ALTER TABLE client_data DROP CONSTRAINT fk_client_user;
ALTER TABLE freelancer_data DROP CONSTRAINT fk_freelancer_user;

ALTER TABLE freelancer_skills 
DROP CONSTRAINT fk_freelancer, 
DROP CONSTRAINT fk_skill;

ALTER TABLE jobs
  DROP CONSTRAINT fk_job_client,
  DROP CONSTRAINT fk_job_freelancer;

ALTER TABLE applications
  DROP CONSTRAINT fk_application_user,
  DROP CONSTRAINT fk_application_job;

ALTER TABLE attachments
  DROP CONSTRAINT fk_attachment_application;
