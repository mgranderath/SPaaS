# SPaaS (Small Product as a Service) [WIP]

A lightweight Heroku like PaaS.

# How to deploy

```
docker run -d \ 
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/<spaas-directory>:/root/.spaas \
    -e HOST_CONFIG_FOLDER='~/<spaas-directory>'
    --label traefik.frontend.rule=Host:spaas.<your-domain>.<your-domain-extension> \
    --name spaas mgranderath/spaas
```