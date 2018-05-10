--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.2
-- Dumped by pg_dump version 9.6.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: conversations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE conversations (
    id uuid NOT NULL,
    status character varying(255) DEFAULT 'pending'::character varying NOT NULL,
    user_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE conversations OWNER TO postgres;

--
-- Name: messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE messages (
    id uuid NOT NULL,
    body text DEFAULT ''::text NOT NULL,
    account_sid character varying(255) DEFAULT ''::character varying NOT NULL,
    message_sid character varying(255) DEFAULT ''::character varying NOT NULL,
    messaging_service_sid character varying(255) DEFAULT ''::character varying NOT NULL,
    sms_message_sid character varying(255) DEFAULT ''::character varying NOT NULL,
    sms_sid character varying(255) DEFAULT ''::character varying NOT NULL,
    reciever_no character varying(255) NOT NULL,
    sender_no character varying(255) NOT NULL,
    sender_city character varying(255) NOT NULL,
    sender_country character varying(255) NOT NULL,
    sender_state character varying(255) NOT NULL,
    sender_zip character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    conversation_id uuid NOT NULL
);


ALTER TABLE messages OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE schema_migration (
    version character varying(255) NOT NULL
);


ALTER TABLE schema_migration OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    id uuid NOT NULL,
    user_name character varying(255) NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    phone_no character varying(255) DEFAULT ''::character varying NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: conversations conversations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY conversations
    ADD CONSTRAINT conversations_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: conversations_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX conversations_status_idx ON conversations USING btree (status);


--
-- Name: conversations_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX conversations_user_id_idx ON conversations USING btree (user_id);


--
-- Name: messages_reciever_no_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX messages_reciever_no_idx ON messages USING btree (reciever_no);


--
-- Name: messages_sender_no_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX messages_sender_no_idx ON messages USING btree (sender_no);


--
-- Name: users_user_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_user_name_idx ON users USING btree (user_name);


--
-- Name: version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX version_idx ON schema_migration USING btree (version);


--
-- Name: conversations conversations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY conversations
    ADD CONSTRAINT conversations_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


--
-- Name: messages messages_conversations_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY messages
    ADD CONSTRAINT messages_conversations_id_fk FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

