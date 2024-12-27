This is the backend component of Ludivault.

For the deployment files, see [this repo](https://github.com/KowalskiPiotr98/ludivault-deploy).

> [!NOTE]
> If you want to read more about Ludivault or check out the frontend source code, please go to [this repository](https://github.com/KowalskiPiotr98/ludivault-web) instead.

## Configuration
The following values can be used to configure the API:
- `GIN_MODE` - please refer to Gin documentation; when in doubt set to `release`,
- `LUDIVAULT_DB` - connection string for the database, more details available in the [Postgres docs](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING); you can use `"host=postgres user=ludivault dbname=ludivault password=ludivault sslmode=disable"` as inspiration (just remember to change the password),
- `LUDIVAULT_LISTEN` - defines an interface at which the application listens for requests; defaults to `localhost:5500` if not set.
