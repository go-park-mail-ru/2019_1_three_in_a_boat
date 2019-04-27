-- Creates database, creates role hexagon, grants permissions, creates tables
--
DO
    $$
        DECLARE
            pwd         TEXT := 'ChangeBeforeDeploying';
            create_role TEXT :=
                concat('CREATE ROLE hexagon WITH LOGIN
       PASSWORD ''', pwd,
                       ''' CONNECTION LIMIT -1');
        BEGIN
            IF pwd = 'ChangeBeforeDeploying'::TEXT THEN -- dont change this line lmao
                RAISE WARNING 'You are using the debugging password!';
            END IF;
            EXECUTE create_role;
        END
        $$;

CREATE DATABASE hexagon WITH OWNER = postgres;

\c hexagon;

CREATE TABLE message (
  id BIGSERIAL PRIMARY KEY,
  uid INT, message TEXT,
  created TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

GRANT ALL PRIVILEGES ON DATABASE hexagon TO hexagon;
GRANT USAGE ON SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO hexagon;
GRANT ALL PRIVILEGES ON TABLE message TO hexagon;