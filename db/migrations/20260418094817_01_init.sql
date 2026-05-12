-- migrate:up

CREATE OR REPLACE FUNCTION set_updated()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TYPE auth.role AS ENUM ('superuser', 'staff', 'manager', 'analist');

CREATE TABLE auth.user (
  id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username      TEXT        NOT NULL UNIQUE,
  password      TEXT        NOT NULL,
  email         TEXT        NOT NULL UNIQUE,
  role          auth.role   NOT NULL DEFAULT 'analist',
  created       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_logged_in TIMESTAMPTZ
);

CREATE TABLE auth.session (
  id       BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id  BIGINT      NOT NULL REFERENCES auth.user(id) ON DELETE CASCADE,
  key      TEXT        NOT NULL UNIQUE,
  created  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires  TIMESTAMPTZ NOT NULL
);

CREATE TRIGGER trg_user_updated
  BEFORE UPDATE ON auth.user
  FOR EACH ROW EXECUTE FUNCTION set_updated();

CREATE TABLE auth.apikey (
  hashed_key  TEXT         PRIMARY KEY,
  user_id     BIGINT       NOT NULL REFERENCES auth.user(id) ON DELETE CASCADE,
  name        TEXT         NOT NULL,
  active      BOOLEAN      NOT NULL DEFAULT TRUE,
  created     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  expiry      TIMESTAMPTZ,
  revoked     BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE SCHEMA IF NOT EXISTS org;

CREATE TABLE org.company (
  id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name          TEXT        NOT NULL UNIQUE,
  display_name  TEXT        NOT NULL,
  created       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE org.analist (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id     BIGINT      NOT NULL REFERENCES auth.user(id) ON DELETE CASCADE,
  first_name  TEXT        NOT NULL,
  last_name   TEXT        NOT NULL,
  company_id  BIGINT      NOT NULL REFERENCES org.company(id) ON DELETE CASCADE,
  created     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE SCHEMA IF NOT EXISTS ngin;

CREATE TABLE ngin.storage (
  id       BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name     TEXT        NOT NULL UNIQUE,
  uri      TEXT        NOT NULL UNIQUE,
  created  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_storage_updated
  BEFORE UPDATE ON ngin.storage
  FOR EACH ROW EXECUTE FUNCTION set_updated();

CREATE TABLE ngin.storage_analist (
  storage_id  BIGINT NOT NULL REFERENCES ngin.storage(id)  ON DELETE CASCADE,
  analist_id  BIGINT NOT NULL REFERENCES org.analist(id)   ON DELETE CASCADE,
  PRIMARY KEY (storage_id, analist_id)
);

CREATE TABLE ngin.notebook (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  analist_id  BIGINT      NOT NULL REFERENCES org.analist(id),
  title       TEXT,
  chat        JSONB       NOT NULL DEFAULT '[]',
  created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_notebook_updated
  BEFORE UPDATE ON ngin.notebook
  FOR EACH ROW EXECUTE FUNCTION set_updated();

CREATE TYPE ngin.model_provider AS ENUM ('xai', 'anthropic', 'openai');

CREATE TABLE ngin.model (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name        TEXT                NOT NULL,
  provider    ngin.model_provider NOT NULL,
  api_key     TEXT                NOT NULL,
  company_id  BIGINT              NOT NULL REFERENCES org.company(id) ON DELETE CASCADE,
  created     TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

-- migrate:down

DROP TRIGGER IF EXISTS trg_notebook_updated ON ngin.notebook;
DROP TRIGGER IF EXISTS trg_storage_updated  ON ngin.storage;

DROP TABLE IF EXISTS ngin.model;
DROP TABLE IF EXISTS ngin.storage_analist;
DROP TABLE IF EXISTS ngin.notebook;
DROP TABLE IF EXISTS ngin.storage;
DROP TYPE  IF EXISTS ngin.model_provider;
DROP SCHEMA IF EXISTS ngin;

DROP TABLE IF EXISTS org.analist;
DROP TABLE IF EXISTS org.company;
DROP SCHEMA IF EXISTS org;

DROP TRIGGER IF EXISTS trg_user_updated ON auth.user;
DROP TABLE IF EXISTS auth.apikey;
DROP TABLE IF EXISTS auth.session;
DROP TABLE IF EXISTS auth.user;
DROP TYPE  IF EXISTS auth.role;
DROP SCHEMA IF EXISTS auth;
DROP FUNCTION IF EXISTS set_updated;
DROP EXTENSION IF EXISTS pgcrypto;
