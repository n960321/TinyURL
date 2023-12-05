CREATE SEQUENCE url_infos_id_seq
    INCREMENT BY 1
	MINVALUE 1000000000
	MAXVALUE 9223372036854775807
	START 1000000000
	CACHE 1
	NO CYCLE;

CREATE TABLE url_infos (
    id integer NOT NULL DEFAULT nextval('url_infos_id_seq'),
    url varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamp,
    CONSTRAINT url_infos_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE url_infos_id_seq OWNED BY url_infos.id;
