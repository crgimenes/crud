CREATE TABLE client_contacts
(
    id serial NOT NULL,
    name character varying(200) NOT NULL,
    fone character varying(100),
    email character varying(200),
    client_id integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT client_id_fk FOREIGN KEY (client_id)
        REFERENCES public.client (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER TABLE client_contacts
    OWNER to postgres;