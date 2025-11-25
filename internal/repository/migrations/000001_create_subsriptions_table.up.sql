CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SEQUENCE IF NOT EXISTS public.subscriptions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.subscriptions (
    id integer NOT NULL DEFAULT nextval('public.subscriptions_id_seq'::regclass),
    service_name text NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    start_date text NOT NULL,
    end_date text 
);
