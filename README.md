# SPaaS (Small Product as a Service) [WIP]

A lightweight Heroku like PaaS.

# How to deploy

```
docker run -d \ 
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.spaas-server:/root/.spaas \
    --label traefik.frontend.rule=Host:spaas.<your-domain>.<your-domain-extension> \
    --name spaas mgranderath/spaas
```