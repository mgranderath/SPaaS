<img src="docs/logo.svg" width="400">

# SPaaS (Small Product as a Service)

A lightweight Heroku like PaaS.

Link to [CLI](https://github.com/mgranderath/SPaaS-cli)

# Info

Currently only nodejs deployments 

# Install Server

### Prerequisites
- Docker
- SSH into Server

## With Domain

You have to have your domain setup already.

Execute on Server:
```
docker run -d \ 
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/<spaas-directory>:/root/.spaas \
    -e HOST_CONFIG_FOLDER='~/<spaas-directory>'
    --label traefik.frontend.rule=Host:spaas.<your-domain>.<your-domain-extension> \
    --name spaas mgranderath/spaas
```

Now you should be able to access the dashboard at `http://spaas.<your-domain>.<your-domain-extension>`

## Without Domain

Execute on Server:
```
docker run -d \ 
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/<spaas-directory>:/root/.spaas \
    -e HOST_CONFIG_FOLDER='~/<spaas-directory>'
    --label traefik.frontend.rule=PathPrefixStrip:/spaas \
    --name spaas mgranderath/spaas
```
Now you should be able to access the dashboard at `http://<your-ip>/spaas` 

# Next Step (Modifying config)

Next you should modify the config to your liking. There should be one created in `~/<spaas-directory>`. It's nane is `.spaas.json`.

You can modify the following settings:

| Setting | Default | Explained |
| ------- | ------- | --------- |
| `username` | spaas | this is the username you use to login |
| `letsencrypt` | false | use automatic letsencrypt ssl certs (only with domains) | 
| `letsEncryptEmail` | "example@example.com" | this is the email that will be used to request certs |
| `domain` | example.com | the domain to be used for deployments |
| `useDomain` | false | use the domain for deployments |

If you change the `letsencrypt` or `letsEncryptEmail` settings you should remove the traefik container (`docker rm -f spaas-traefik`) and then restart SPaaS (`docker restart spaas`).

## Important Info

- The default login credentials are `spaas`:`smallpaas`
- Without a domain the deployed apps can be accessed under `<your-ip>/spaas/<app-name>`
- With a domain the deployed apps can be accessed under `<app-name>.<your-domain>.<your-domain-extension>`
- **!!! You should change the password as soon as possible. Either using the dashboard or the cli application !!!**

# Deploying a App

**!!! Only nodejs is supported for now !!!**

## Create a new app

When you create a new app using either the cli or the dashboard you should get back a path called `RepoPath`. This is the location of the git repo on the server. 

You should add the following as a git remote for the project you want to deploy:
```bash
git remote add spaas <ssh-server-username>@<server>:<repo-path>
```

## Before Deploying

You have to specify the start command for the server to know which command to run to start the web server.

Add a file `spaas.json` in the project directory.

```json
{
  "start": "<app-start-command>"
}
```

## Deploying

To deploy the app you just have to follow the usual procedure to push to a git repository.

1. `git add .`
2. `git commit -m "message"`
3. `git push spaas master`

If the push succeeds and you do not see any error you should be able to see your app on the address that is specified above depending on your configuration.