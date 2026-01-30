Using the `efactura-cli` and `etransport-cli` commands to get an ANAF Auth Token
=================================

This guide should help you generate an access token for your local machine
(great for testing and CLI usage), but with small adjustments can be used for
production too.

This guide assumes that you have a valid and active qualified certificate on a
token device from a company accepted by ANAF and that the token devices is
supported by your browser (this is usually best supported on Windows, although
I use a DigiSign Token with [SafeNet Authentication Client 10.9](https://www.digicert.com/StaticFiles/Linux_SAC_10.9_GA.zip)
for Linux it's working great so far).

The `auth` subcommand is identical for both `efactura-cli` and `etransport-cli`,
so for the purposes of this guide, the two commands are completely interchangable
(as long as both of them are installed from the same version / git commit hash).

1. ANAF OAuth Apps registered in the ANAF portal need a callback URL - a valid
   HTTPS URL where you will be redirected during OAuth2
   [Authorization Code flow](https://auth0.com/docs/get-started/authentication-and-authorization-flow/authorization-code-flow)
   after authenticating with your qualified certificate and where to OAuth
   Authorization code is sent (the code needs to be exchanged for an access token).
   If you are using a production setup or you have a DNS domain that you have
   acccess to, you can go directly to step 3.

> [!IMPORTANT]
> Once you created an OAuth App, you cannot edit it, so you cannot update the
> callback URL, but you can provide multiple callback URLs for an app.

2. Create a [LocalTunnel](https://theboroer.github.io/localtunnel-www/) -
   basically expose your machine to the internet so you can receive the
   Authorization code:

```
# Install the localtun command
go install github.com/printesoi/go-localtunnel/cmd/localtun@v1.0.0

# Start a localtunnel by using a custom hard to guess subdomain:
localtun --host localhost --port 5000 --subdomain "$(uuigen -r)"
```

You can see your custom publicly accessible domain in the log output from the
localtun command:

```
2026/01/28 00:29:36 registering tunnel: https://localtunnel.me/528e730b-16c5-42fe-85c7-7e707f9ebb18
2026/01/28 00:29:36 registered tunnel: https://528e730b-16c5-42fe-85c7-7e707f9ebb18.loca.lt
...
```

in the command above the full domain is: `https://528e730b-16c5-42fe-85c7-7e707f9ebb18.loca.lt`.
TL;DR All requests to `https://528e730b-16c5-42fe-85c7-7e707f9ebb18.loca.lt`
will te reverse proxied to your `localhost` on port `5000`.

> [!IMPORTANT]
> `528e730b-16c5-42fe-85c7-7e707f9ebb18` is just an example value. For security
> reasons, I recommend actually using a hard to guess value (like the output of `uuigen -r`)
> and replace the value in all commands bellow.

3. Register a test OAuth app in the ANAF "SPV > Editare profil Oauth" portal:
   - Use whatever you want for the App name ("Denumire aplicație").
   - For the Callback URL use the LocalTunnel custom domain with the callback
     path: `https://528e730b-16c5-42fe-85c7-7e707f9ebb18.loca.lt/callback`
   - For the service ("Serviciu") select whatever you need: `E-Factura`,
     `E-Transport` or both.
   - Click on "Generare Client ID". This will create a new app the you should
     be able to see in the "Client ID-uri Oauth existente pentru contul dvs."
     list.
   - You need to copy and store the `Client ID` and `Client Secret` for your
     app.

> [!IMPORTANT]
> The `Client ID` and `Client Secret` are secrets and should be treated as such.
> Be very careful where and how you store them, don't commit them to your repo
> and don't share them, unless you know what you're doing.

4. Run the `efactura-cli`/`etransport-cli` authorize callback handler.
   Use the client ID and client secret for the ANAF OAuth app either via
   command line flags (`--oauth-client-id`, `--oauth-client-secret`) or via a
   YAML config using `--config` (each flag is automatically mapped to a env
   variable (`--oauth-client-id` is mapped to `EFACTURA_OAUTH_CLIENT_ID` env variable)
   or to a config variable (`--oauth-client-id` is mapped to `oauth-client-id`
   in YAML) and `--oauth-redirect-url` (the callback URL):

```
efactura-cli auth [...] --addr localhost:5001 authorize-server --callback-path /callback
```

Replace `localhost:5000` with the host and port you use for the `localtun`
command, and `/callback` with the path you used in the callback URL when you
created the ANAF OAuth app.

5. Get an Authorize URL (it needs the OAuth client ID, secret and redirect URL
   and you can provide them via one of the ways decribed above):

```
efactura-cli auth get-authorize-link [...]
```

For now, this command does not automatically open your browser (for various
reasons). You need to copy / click (depending on how advanced you terminal is)
the URL and open in your browser. Before you open this URL in your browser,
ensure that your token device that stores your qualified certificate
is plugged in and working. The link should look like `https://logincert.anaf.ro/anaf-oauth2/v1/authorize?client_id=...`
and will redirect to a page where you would authenticate using your cerficate
and will redirect to your configured to your app callback URL, which (if you
configured everything correctly) should redirect to the authorize handler that
we ran on step 4 and pass a `code` query param. The handler will automatically
parse the passed code and exchange for an access token and will print the OAuth
access token.

If, for any reason this did not work, but you have a valid authorization code,
you change manually exchange it for an access token using:
```
efactura-cli auth exchange-code [...] --oauth-device-code $code
```

6. After you've generated an access token (it should look like `{"access_token":....`),
   you can either pass to any subcommand via `--access-token` CLI flag, the
   `EFACTURA_ACCESS_TOKEN` env variable or `access-token` config var.
