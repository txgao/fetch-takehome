--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Homebrew)
-- Dumped by pg_dump version 14.13 (Homebrew)

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: item; Type: TABLE; Schema: public; Owner: txgao
--

CREATE TABLE public.item (
    item_uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    price numeric NOT NULL,
    short_description character varying(255)
);


ALTER TABLE public.item OWNER TO txgao;

--
-- Name: receipt; Type: TABLE; Schema: public; Owner: txgao
--

CREATE TABLE public.receipt (
    receipt_uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    total numeric NOT NULL,
    purchase_date date NOT NULL,
    purchase_time time without time zone NOT NULL,
    retailer character varying(255) NOT NULL
);


ALTER TABLE public.receipt OWNER TO txgao;

--
-- Name: receipt_items; Type: TABLE; Schema: public; Owner: txgao
--

CREATE TABLE public.receipt_items (
    item_uuid uuid NOT NULL,
    receipt_uuid uuid NOT NULL
);


ALTER TABLE public.receipt_items OWNER TO txgao;

--
-- Name: item item_pkey; Type: CONSTRAINT; Schema: public; Owner: txgao
--

ALTER TABLE ONLY public.item
    ADD CONSTRAINT item_pkey PRIMARY KEY (item_uuid);


--
-- Name: receipt_items receipt_items_pkey; Type: CONSTRAINT; Schema: public; Owner: txgao
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_pkey PRIMARY KEY (item_uuid, receipt_uuid);


--
-- Name: receipt receipt_pkey; Type: CONSTRAINT; Schema: public; Owner: txgao
--

ALTER TABLE ONLY public.receipt
    ADD CONSTRAINT receipt_pkey PRIMARY KEY (receipt_uuid);


--
-- Name: receipt_items receipt_items_item_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: txgao
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_item_uuid_fkey FOREIGN KEY (item_uuid) REFERENCES public.item(item_uuid);


--
-- Name: receipt_items receipt_items_receipt_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: txgao
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_receipt_uuid_fkey FOREIGN KEY (receipt_uuid) REFERENCES public.receipt(receipt_uuid);


--
-- PostgreSQL database dump complete
--

