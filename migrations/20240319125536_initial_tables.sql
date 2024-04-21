-- migrate:up

-- create app user
DO
$do$
BEGIN
	IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'sd_appuser') THEN
		CREATE USER sd_appuser WITH ENCRYPTED PASSWORD 'splitd';
	END IF;
END
$do$;

-- update access
REVOKE ALL ON DATABASE splitdecision FROM public;

GRANT CONNECT
ON DATABASE splitdecision
TO sd_appuser;

-- GRANT USAGE
-- ON SCHEMA public
-- TO sd_appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE, TRIGGER
ON TABLES
TO sd_appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT EXECUTE
ON FUNCTIONS
TO sd_appuser;

-- create functions
CREATE OR REPLACE FUNCTION manage_table_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$ language 'plpgsql';

-- create tables
CREATE TABLE IF NOT EXISTS auth_types(
	auth_type VARCHAR(64) PRIMARY KEY
);


CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	email VARCHAR(320) NOT NULL,
	auth_type VARCHAR(64) NOT NULL REFERENCES auth_types(auth_type),
	token VARCHAR(256) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();

CREATE TABLE IF NOT EXISTS sessions(
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id),
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON sessions FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();

CREATE TABLE IF NOT EXISTS decisions(
	id UUID PRIMARY KEY,
	owner_id UUID NOT NULL REFERENCES users(id),
	public BOOLEAN NOT NULL DEFAULT false,
	prompt VARCHAR(512) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON decisions FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS choices(
	id UUID PRIMARY KEY,
	owner_id UUID DEFAULT NULL REFERENCES users(id),
	name VARCHAR(128) NOT NULL,
	image_url TEXT DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON choices FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();

CREATE TABLE IF NOT EXISTS brackets(
	id UUID PRIMARY KEY,
	owner_id UUID NOT NULL REFERENCES users(id),
	decision_id UUID NOT NULL REFERENCES decisions(id),
	winner_id UUID DEFAULT NULL REFERENCES choices(id),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON brackets FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS bracket_choices(
	bracket_id UUID NOT NULL REFERENCES brackets(id),
	item_id UUID NOT NULL REFERENCES choices(id),
	UNIQUE(bracket_id, item_id)
);


CREATE TABLE IF NOT EXISTS matches(
	id UUID PRIMARY KEY,
	bracket_id UUID NOT NULL REFERENCES brackets(id),
	left_choice_id UUID NOT NULL REFERENCES choices(id),
	right_choice_id UUID NOT NULL REFERENCES choices(id),
	winner_id UUID DEFAULT NULL REFERENCES choices(id),
	round INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON matches FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


-- seed initial data
INSERT INTO auth_types
	(auth_type)
VALUES
	('basic'),
	('google')
ON CONFLICT DO NOTHING;


-- migrate:down
DROP TRIGGER IF EXISTS manage_updated_at ON users;
DROP TRIGGER IF EXISTS manage_updated_at ON brackets;
DROP TRIGGER IF EXISTS manage_updated_at ON decisions;
DROP TRIGGER IF EXISTS manage_updated_at ON choices;
DROP TRIGGER IF EXISTS manage_updated_at ON matches;

DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS bracket_choices;
DROP TABLE IF EXISTS brackets;
DROP TABLE IF EXISTS choices;
DROP TABLE IF EXISTS decisions;
DROP TABLE IF EXISTS users;

REASSIGN OWNED BY sd_appuser TO root;
DROP OWNED BY sd_appuser;
DROP USER IF EXISTS sd_appuser;
