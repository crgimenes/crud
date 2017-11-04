# gocrud
Simple master detail web form with Golang

- Postgresql --> Database

- pREST --> API REST

- App frontend --> serv application pages (html)

- Caddy --> Proxy


# Install

...

# Start Postgresql server

...

# Run database migrations

```
export PREST_HTTP_PORT=3000
export PREST_PG_HOST=127.0.0.1
export PREST_PG_USER=postgres
#export PREST_PG_PASS=
export PREST_PG_DATABASE=gocrud
export PREST_PG_PORT=5432
export PREST_JWT_KEY="your jwt key here"
prest migrate up
```

# Start pREST server

```
prest
```

# Start frontend app



# Start Caddy

```
caddy -conf ./caddy/Caddyfile
```