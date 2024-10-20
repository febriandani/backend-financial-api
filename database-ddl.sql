-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id bigserial NOT NULL,
	"userId" int4 NOT NULL,
	"name" varchar(100) NOT NULL,
	username varchar(20) NULL,
	email varchar(100) NULL,
	phone_number varchar(15) NULL,
	"password" varchar(100) NULL,
	created_at timestamp NULL,
	created_by varchar(100) DEFAULT 'system'::character varying NULL,
	updated_at timestamp NULL,
	updated_by varchar(100) DEFAULT 'system'::character varying NULL,
	is_active bool DEFAULT true NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- public.transactions definition

-- Drop table

-- DROP TABLE public.transactions;

CREATE TABLE public.transactions (
	id bigserial NOT NULL,
	user_id int4 NOT NULL,
	category_id int4 NOT NULL,
	category_type varchar(10) NOT NULL,
	amount varchar DEFAULT '0'::character varying NULL,
	description varchar NOT NULL,
	created_at timestamp NULL,
	created_by varchar(100) DEFAULT 'system'::character varying NULL,
	updated_at timestamp NULL,
	updated_by varchar(100) DEFAULT 'system'::character varying NULL,
	CONSTRAINT transactions_pkey PRIMARY KEY (id)
);

-- public.categories definition

-- Drop table

-- DROP TABLE public.categories;

CREATE TABLE public.categories (
	id bigserial NOT NULL,
	user_id int4 NOT NULL,
	category_type varchar(10) NOT NULL,
	category_name varchar(100) NOT NULL,
	category_description varchar(100) NOT NULL,
	created_at timestamp NULL,
	created_by varchar(100) DEFAULT 'system'::character varying NULL,
	updated_at timestamp NULL,
	updated_by varchar(100) DEFAULT 'system'::character varying NULL,
	CONSTRAINT categories_pkey PRIMARY KEY (id)
);