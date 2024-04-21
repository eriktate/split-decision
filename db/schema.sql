SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: manage_table_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.manage_table_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auth_types; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.auth_types (
    auth_type character varying(64) NOT NULL
);


--
-- Name: bracket_choices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bracket_choices (
    bracket_id uuid NOT NULL,
    item_id uuid NOT NULL
);


--
-- Name: brackets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.brackets (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    decision_id uuid NOT NULL,
    winner_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: choices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.choices (
    id uuid NOT NULL,
    owner_id uuid,
    name character varying(128) NOT NULL,
    image_url text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: decisions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.decisions (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    public boolean DEFAULT false NOT NULL,
    prompt character varying(512) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--
-- Name: matches; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.matches (
    id uuid NOT NULL,
    bracket_id uuid NOT NULL,
    left_choice_id uuid NOT NULL,
    right_choice_id uuid NOT NULL,
    winner_id uuid,
    round integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(320) NOT NULL,
    auth_type character varying(64) NOT NULL,
    token character varying(256) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: auth_types auth_types_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth_types
    ADD CONSTRAINT auth_types_pkey PRIMARY KEY (auth_type);


--
-- Name: bracket_choices bracket_choices_bracket_id_item_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bracket_choices
    ADD CONSTRAINT bracket_choices_bracket_id_item_id_key UNIQUE (bracket_id, item_id);


--
-- Name: brackets brackets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.brackets
    ADD CONSTRAINT brackets_pkey PRIMARY KEY (id);


--
-- Name: choices choices_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.choices
    ADD CONSTRAINT choices_pkey PRIMARY KEY (id);


--
-- Name: decisions decisions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.decisions
    ADD CONSTRAINT decisions_pkey PRIMARY KEY (id);


--
-- Name: matches matches_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.matches
    ADD CONSTRAINT matches_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: brackets manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.brackets FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: choices manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.choices FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: decisions manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.decisions FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: matches manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.matches FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: sessions manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.sessions FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: users manage_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER manage_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.manage_table_updated_at();


--
-- Name: bracket_choices bracket_choices_bracket_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bracket_choices
    ADD CONSTRAINT bracket_choices_bracket_id_fkey FOREIGN KEY (bracket_id) REFERENCES public.brackets(id);


--
-- Name: bracket_choices bracket_choices_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bracket_choices
    ADD CONSTRAINT bracket_choices_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.choices(id);


--
-- Name: brackets brackets_decision_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.brackets
    ADD CONSTRAINT brackets_decision_id_fkey FOREIGN KEY (decision_id) REFERENCES public.decisions(id);


--
-- Name: brackets brackets_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.brackets
    ADD CONSTRAINT brackets_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: brackets brackets_winner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.brackets
    ADD CONSTRAINT brackets_winner_id_fkey FOREIGN KEY (winner_id) REFERENCES public.choices(id);


--
-- Name: choices choices_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.choices
    ADD CONSTRAINT choices_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: decisions decisions_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.decisions
    ADD CONSTRAINT decisions_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: matches matches_bracket_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.matches
    ADD CONSTRAINT matches_bracket_id_fkey FOREIGN KEY (bracket_id) REFERENCES public.brackets(id);


--
-- Name: matches matches_left_choice_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.matches
    ADD CONSTRAINT matches_left_choice_id_fkey FOREIGN KEY (left_choice_id) REFERENCES public.choices(id);


--
-- Name: matches matches_right_choice_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.matches
    ADD CONSTRAINT matches_right_choice_id_fkey FOREIGN KEY (right_choice_id) REFERENCES public.choices(id);


--
-- Name: matches matches_winner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.matches
    ADD CONSTRAINT matches_winner_id_fkey FOREIGN KEY (winner_id) REFERENCES public.choices(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users users_auth_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_auth_type_fkey FOREIGN KEY (auth_type) REFERENCES public.auth_types(auth_type);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240319125536');
