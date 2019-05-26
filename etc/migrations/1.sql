CREATE TABLE message (
  id BIGSERIAL PRIMARY KEY,
  uid INT REFERENCES account(uid),
  message TEXT,
  created TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);