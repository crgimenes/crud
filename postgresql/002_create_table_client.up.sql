CREATE TABLE client
(
    id serial NOT NULL,
    name character varying(150) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE client
    OWNER to postgres;