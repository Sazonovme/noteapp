ALTER TABLE users 
    RENAME COLUMN login TO email;

ALTER TABLE groups 
    RENAME COLUMN user_login TO user_email;

ALTER TABLE notes 
    RENAME COLUMN user_login TO user_email;

ALTER TABLE refreshsessions 
    RENAME COLUMN user_login TO user_email;