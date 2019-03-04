CREATE ROLE hexagon WITH LOGIN 
                         PASSWORD 'ChangeBeforeDeploying' -- imma forget to change it 100% of the time lmfao
                         CONNECTION LIMIT -1;

CREATE DATABASE hexagon WITH OWNER = postgres;

\c hexagon;

CREATE TYPE gender_t AS ENUM ('male', 'female', 'other');

CREATE TABLE IF NOT EXISTS account (
  uid         SERIAL                      PRIMARY KEY,
  username    VARCHAR(32)                 NOT NULL CHECK (username <> ''),
  email       VARCHAR(254)                NOT NULL CHECK (email <> '')
);

CREATE TABLE IF NOT EXISTS profile (
  uid         integer                     PRIMARY KEY REFERENCES account,
  first_name  VARCHAR(32),
  last_name   VARCHAR(32),
  high_score  int                         CHECK (high_score >= 0),
  gender      gender_t,
  userpic     VARCHAR(64),
  birth_data  date,
  signup_date timestamp(0) with time zone NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS author (
  uid         integer                     PRIMARY KEY REFERENCES account,
  devInfo     varchar(128),
  img         varchar(128),
  description text
);

GRANT ALL PRIVILEGES ON DATABASE hexagon TO hexagon;
GRANT USAGE ON SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO hexagon;