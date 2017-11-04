# gocrud
Simple master detail web form with Golang

- Postgresql --> Database

- pREST --> API REST

- gatekeeper --> Autentication server

- App frontend --> serv application pages (html)

- Caddy --> Proxy


# Install

## pREST

```
go get github.com/prest/prest
```

## Caddy

```
go get github.com/mholt/caddy/caddy
```

# Start Postgresql server

- plataform dependent

## Create database

```
createdb gocrud
```

# Run database migrations

attention: To run pREST you need a precreated database

```
export PREST_HTTP_PORT=3000
export PREST_PG_HOST=127.0.0.1
export PREST_PG_USER=postgres
#export PREST_PG_PASS=
export PREST_PG_DATABASE=gocrud
export PREST_PG_PORT=5432
export PREST_JWT_KEY="your jwt key here"
export PREST_MIGRATIONS=./postgresql
export PREST_QUERIES_LOCATION=./queries
prest migrate up
```

# Start pREST server

```
prest
```

# Start Gatekeeper

```
gatekeeper
```

## Teste Gatekeeper

```
curl localhost:4000 
```

Then get the jwt token and put the value on $JWT_TOKEN environment variable, to test pREST.

## Test pREST with authentication

```
curl -i -H 'authorization: Bearer $JWT_TOKEN localhost:3000/databases
```

You must see your databases now.


# Start frontend app



# Start Caddy

```
caddy -conf ./caddy/Caddyfile
```