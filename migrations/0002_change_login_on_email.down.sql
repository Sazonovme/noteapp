ALTER TABLE users 
    RENAME COLUMN email TO login;

ALTER TABLE groups 
    RENAME COLUMN user_email TO user_login;

ALTER TABLE notes 
    RENAME COLUMN user_email TO user_login;

ALTER TABLE refreshsessions 
    RENAME COLUMN user_email TO user_login;
