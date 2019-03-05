-- Creates database, creates role hexagon, grants permissions, creates tables
--
\set pwd 'ChangeBeforeDeploying';
DO
  $$
    BEGIN
      IF :pwd == 'ChangeBeforeDeploying' THEN -- dont change this line lmao
        RAISE WARNING 'You are using the debugging password!';
      END IF;
    END
  $$;

CREATE ROLE hexagon WITH LOGIN
  PASSWORD :pwd
  CONNECTION LIMIT -1;

CREATE DATABASE hexagon WITH OWNER = postgres;

\c hexagon;

CREATE TYPE GENDER_T AS ENUM ('male', 'female', 'other');

CREATE TABLE IF NOT EXISTS account
(
  uid      SERIAL       PRIMARY KEY,
  username VARCHAR(32)  UNIQUE NOT NULL CHECK (username <> ''),
  password BYTEA        NOT NULL CHECK (octet_length(password) <> 0),
  email    VARCHAR(254) UNIQUE NOT NULL CHECK (email <> '')
);

CREATE TABLE IF NOT EXISTS profile
(
  uid         INTEGER                     PRIMARY KEY REFERENCES account
                                          DEFERRABLE INITIALLY DEFERRED,
  first_name  VARCHAR(32),
  last_name   VARCHAR(32),
  high_score  INT                         CHECK (high_score >= 0),
  gender      GENDER_T,
  img         VARCHAR(64),
  birth_date  DATE,
  signup_date TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS author
(
  uid         INTEGER      PRIMARY KEY REFERENCES account DEFERRABLE INITIALLY IMMEDIATE,
  dev_info    VARCHAR(128),
  description TEXT
);

GRANT ALL PRIVILEGES ON DATABASE hexagon TO hexagon;
GRANT USAGE ON SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO hexagon;