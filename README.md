This is the backend component of Ludivault.

For the deployment files, see [this repo](https://github.com/KowalskiPiotr98/ludivault-deploy).

> [!NOTE]
> If you want to read more about Ludivault or check out the frontend source code, please go to [this repository](https://github.com/KowalskiPiotr98/ludivault-web) instead.

## Configuration
The following values can be used to configure the API:
- `GIN_MODE` - please refer to Gin documentation; when in doubt set to `release`,
- `LUDIVAULT_DB` - connection string for the database, more details available in the [Postgres docs](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING); you can use `"host=postgres user=ludivault dbname=ludivault password=ludivault sslmode=disable"` as inspiration (just remember to change the password),
- `LUDIVAULT_LISTEN` - defines an interface at which the application listens for requests; defaults to `localhost:5500` if not set.
- `LUDIVAULT_BASE_ADDRESS` - base public address by which the user will access Ludivault. Used for SSO callback config - does not affect listen address. (example: `https://ludivault.localdomain/`)
- `LUDIVAULT_SESSION_KEY` - secret key used for session tokens encryption. You **MUST** set this to a random, secret value. You can change this value to log out all users at once (requires restart of the application).

### Login providers
The application does not support login with a local account.
All users must log in using an existing external login provider.
Login providers are configured by setting environment variables.
At least one login provider must be configured.

> [!NOTE]
> Once a user logs in using an identity provider, it's not possible to change it (at least for now).
> For that reason, removing an existing login provider is currently not supported without first manually transferring users in some way.

The following login providers are supported (expand for configuration details):

<details>
<summary>Gitea</summary>

- `LUDIVAULT_SSO_GITEA_CLIENT_ID`
- `LUDIVAULT_SSO_GITEA_CLIENT_SECRET`
- `LUDIVAULT_SSO_GITEA_URL` - base url of the Gitea instance

The Gitea is the suggested way of authenticating users for development or test purposes, as it's extremely easy to deploy and use.

Note that only self-hosted Gitea is supported.
</details>

<details>
<summary>Twitch</summary>

- `LUDIVAULT_SSO_TWITCH_CLIENT_ID`
- `LUDIVAULT_SSO_TWITCH_CLIENT_SECRET`
</details>

<details>
<summary>GitHub</summary>

- `LUDIVAULT_SSO_GITHUB_CLIENT_ID`
- `LUDIVAULT_SSO_GITHUB_CLIENT_SECRET`

Note that only the GitHub service is supported.
Self-hosted GitHub Enterprise is not supported (because why would it).
</details>

<details>
<summary>Discord</summary>

- `LUDIVAULT_SSO_DISCORD_CLIENT_ID`
- `LUDIVAULT_SSO_DISCORD_CLIENT_SECRET`
</details>

<details>
<summary>Google</summary>

- `LUDIVAULT_SSO_GOOGLE_CLIENT_ID`
- `LUDIVAULT_SSO_GOOGLE_CLIENT_SECRET`
</details>

<details>
<summary>Steam</summary>

- `LUDIVAULT_SSO_STEAM_CLIENT_SECRET` - API key provided by Steam
</details>

<details>
<summary>Custom OpenID Connect provider</summary>

- `LUDIVAULT_SSO_OIDC_CLIENT_ID`
- `LUDIVAULT_SSO_OIDC_CLIENT_SECRET`
- `LUDIVAULT_SSO_OIDC_DISCOVERY_URL`

This allows you to set up authentication using any provider compatible with the OpenID Connect protocol.

Only one custom provider can be configured at a time.
Changing custom providers is not supported, unless all users and their ids are retained between the providers.
</details>
