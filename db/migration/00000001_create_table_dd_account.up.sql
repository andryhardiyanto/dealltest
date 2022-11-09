-- public.dd_account definition

-- Drop table

-- DROP TABLE public.dd_account;

CREATE TABLE public.dd_account (
	ddacc_id bigserial NOT NULL,
	ddacc_name text NULL DEFAULT ''::text,
	ddacc_email text NULL DEFAULT ''::text,
	ddacc_password text NULL DEFAULT ''::text,
	ddacc_role varchar(20) NULL DEFAULT ''::text,
	created_at timestamptz NULL DEFAULT now(),
	updated_at timestamptz NULL DEFAULT now(),
	CONSTRAINT dd_account_pk PRIMARY KEY (ddacc_id)
);