DROP TABLE IF EXISTS notes;

DROP TABLE IF EXISTS refreshsessions;

DROP TABLE IF EXISTS groups;

DROP TABLE IF EXISTS users;

REVOKE ALL ON schema public FROM notesapp;

REVOKE ALL ON ALL TABLES IN SCHEMA public FROM notesapp;

DROP ROLE notesapp;