\restrict dbmate

-- Dumped from database version 17.9 (Homebrew)
-- Dumped by pg_dump version 17.9 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: auth; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA auth;


--
-- Name: ngin; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA ngin;


--
-- Name: org; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA org;


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: role; Type: TYPE; Schema: auth; Owner: -
--

CREATE TYPE auth.role AS ENUM (
    'superuser',
    'staff',
    'manager',
    'analist'
);


--
-- Name: model_provider; Type: TYPE; Schema: ngin; Owner: -
--

CREATE TYPE ngin.model_provider AS ENUM (
    'xai',
    'anthropic',
    'openai'
);


--
-- Name: set_updated(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated = NOW();
  RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apikey; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE auth.apikey (
    hashed_key text NOT NULL,
    user_id bigint NOT NULL,
    name text NOT NULL,
    active boolean DEFAULT true NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    expiry timestamp with time zone,
    revoked boolean DEFAULT false NOT NULL
);


--
-- Name: session; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE auth.session (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    key text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    expires timestamp with time zone NOT NULL
);


--
-- Name: session_id_seq; Type: SEQUENCE; Schema: auth; Owner: -
--

ALTER TABLE auth.session ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME auth.session_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: user; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE auth."user" (
    id bigint NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    role auth.role DEFAULT 'analist'::auth.role NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    updated timestamp with time zone DEFAULT now() NOT NULL,
    last_logged_in timestamp with time zone
);


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: auth; Owner: -
--

ALTER TABLE auth."user" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME auth.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: model; Type: TABLE; Schema: ngin; Owner: -
--

CREATE TABLE ngin.model (
    id bigint NOT NULL,
    name text NOT NULL,
    provider ngin.model_provider NOT NULL,
    api_key text NOT NULL,
    company_id bigint NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: model_id_seq; Type: SEQUENCE; Schema: ngin; Owner: -
--

ALTER TABLE ngin.model ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME ngin.model_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: notebook; Type: TABLE; Schema: ngin; Owner: -
--

CREATE TABLE ngin.notebook (
    id bigint NOT NULL,
    analist_id bigint NOT NULL,
    title text,
    chat jsonb DEFAULT '[]'::jsonb NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    updated timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: notebook_id_seq; Type: SEQUENCE; Schema: ngin; Owner: -
--

ALTER TABLE ngin.notebook ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME ngin.notebook_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: storage; Type: TABLE; Schema: ngin; Owner: -
--

CREATE TABLE ngin.storage (
    id bigint NOT NULL,
    name text NOT NULL,
    uri text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    updated timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: storage_analist; Type: TABLE; Schema: ngin; Owner: -
--

CREATE TABLE ngin.storage_analist (
    storage_id bigint NOT NULL,
    analist_id bigint NOT NULL
);


--
-- Name: storage_id_seq; Type: SEQUENCE; Schema: ngin; Owner: -
--

ALTER TABLE ngin.storage ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME ngin.storage_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: analist; Type: TABLE; Schema: org; Owner: -
--

CREATE TABLE org.analist (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    company_id bigint NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: analist_id_seq; Type: SEQUENCE; Schema: org; Owner: -
--

ALTER TABLE org.analist ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME org.analist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: company; Type: TABLE; Schema: org; Owner: -
--

CREATE TABLE org.company (
    id bigint NOT NULL,
    name text NOT NULL,
    display_name text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: company_id_seq; Type: SEQUENCE; Schema: org; Owner: -
--

ALTER TABLE org.company ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME org.company_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: apikey apikey_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.apikey
    ADD CONSTRAINT apikey_pkey PRIMARY KEY (hashed_key);


--
-- Name: session session_key_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.session
    ADD CONSTRAINT session_key_key UNIQUE (key);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (id);


--
-- Name: user user_email_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user user_username_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth."user"
    ADD CONSTRAINT user_username_key UNIQUE (username);


--
-- Name: model model_pkey; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.model
    ADD CONSTRAINT model_pkey PRIMARY KEY (id);


--
-- Name: notebook notebook_pkey; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.notebook
    ADD CONSTRAINT notebook_pkey PRIMARY KEY (id);


--
-- Name: storage_analist storage_analist_pkey; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage_analist
    ADD CONSTRAINT storage_analist_pkey PRIMARY KEY (storage_id, analist_id);


--
-- Name: storage storage_name_key; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage
    ADD CONSTRAINT storage_name_key UNIQUE (name);


--
-- Name: storage storage_pkey; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage
    ADD CONSTRAINT storage_pkey PRIMARY KEY (id);


--
-- Name: storage storage_uri_key; Type: CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage
    ADD CONSTRAINT storage_uri_key UNIQUE (uri);


--
-- Name: analist analist_pkey; Type: CONSTRAINT; Schema: org; Owner: -
--

ALTER TABLE ONLY org.analist
    ADD CONSTRAINT analist_pkey PRIMARY KEY (id);


--
-- Name: company company_name_key; Type: CONSTRAINT; Schema: org; Owner: -
--

ALTER TABLE ONLY org.company
    ADD CONSTRAINT company_name_key UNIQUE (name);


--
-- Name: company company_pkey; Type: CONSTRAINT; Schema: org; Owner: -
--

ALTER TABLE ONLY org.company
    ADD CONSTRAINT company_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user trg_user_updated; Type: TRIGGER; Schema: auth; Owner: -
--

CREATE TRIGGER trg_user_updated BEFORE UPDATE ON auth."user" FOR EACH ROW EXECUTE FUNCTION public.set_updated();


--
-- Name: notebook trg_notebook_updated; Type: TRIGGER; Schema: ngin; Owner: -
--

CREATE TRIGGER trg_notebook_updated BEFORE UPDATE ON ngin.notebook FOR EACH ROW EXECUTE FUNCTION public.set_updated();


--
-- Name: storage trg_storage_updated; Type: TRIGGER; Schema: ngin; Owner: -
--

CREATE TRIGGER trg_storage_updated BEFORE UPDATE ON ngin.storage FOR EACH ROW EXECUTE FUNCTION public.set_updated();


--
-- Name: apikey apikey_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.apikey
    ADD CONSTRAINT apikey_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth."user"(id) ON DELETE CASCADE;


--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth."user"(id) ON DELETE CASCADE;


--
-- Name: model model_company_id_fkey; Type: FK CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.model
    ADD CONSTRAINT model_company_id_fkey FOREIGN KEY (company_id) REFERENCES org.company(id) ON DELETE CASCADE;


--
-- Name: notebook notebook_analist_id_fkey; Type: FK CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.notebook
    ADD CONSTRAINT notebook_analist_id_fkey FOREIGN KEY (analist_id) REFERENCES org.analist(id);


--
-- Name: storage_analist storage_analist_analist_id_fkey; Type: FK CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage_analist
    ADD CONSTRAINT storage_analist_analist_id_fkey FOREIGN KEY (analist_id) REFERENCES org.analist(id) ON DELETE CASCADE;


--
-- Name: storage_analist storage_analist_storage_id_fkey; Type: FK CONSTRAINT; Schema: ngin; Owner: -
--

ALTER TABLE ONLY ngin.storage_analist
    ADD CONSTRAINT storage_analist_storage_id_fkey FOREIGN KEY (storage_id) REFERENCES ngin.storage(id) ON DELETE CASCADE;


--
-- Name: analist analist_company_id_fkey; Type: FK CONSTRAINT; Schema: org; Owner: -
--

ALTER TABLE ONLY org.analist
    ADD CONSTRAINT analist_company_id_fkey FOREIGN KEY (company_id) REFERENCES org.company(id) ON DELETE CASCADE;


--
-- Name: analist analist_user_id_fkey; Type: FK CONSTRAINT; Schema: org; Owner: -
--

ALTER TABLE ONLY org.analist
    ADD CONSTRAINT analist_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth."user"(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict dbmate


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20260418094817');
