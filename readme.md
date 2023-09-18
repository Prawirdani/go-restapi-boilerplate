### Personal project boilerplate for go restful api server.
Based on [Chi](https://github.com/go-chi) router.
additional:
- Postgres driver: [pgx](github.com/jackc/pgx/v5)
- Middleware:
    - [Secure](https://github.com/unrolled/secure)
    - [Cors](https://github.com/rs/cors)
- Utility:
    - slog logger `go 1.21`
    - [Request Validator](https://github.com/go-playground/validator/v10)
- Docs:
    - [Swagger](https://github.com/swaggo/http-swagger)

Used most of the time for Layered Architecture, probably suitable too for DDD. I'm not going to write any documentation, so if you want to use this, go play around and find out! 