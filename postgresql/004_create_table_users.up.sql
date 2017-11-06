CREATE TABLE users
(
    id serial NOT NULL,
    username character varying(200) NOT NULL,
    password text NOT NULL,
    customer_id integer NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE public.users
    OWNER to postgres;

-- default username=gocrud@example.com password=1234
INSERT INTO public.users(
	id, username, password, customer_id)
	VALUES (
    1, 
    'gocrud@example.com', 
    'BB06E371F6B25E45A7F2082115346F5EB78E4501BDFECD59F8DBC4643A7355BA', 
    1);