# docker-farmer

Custom build tool to remove docker containers and databases from services.

## Config

Example configuration:

```js
{
    "domain": "test.example.com",
    "docker": {
        "host": "unix:///var/run/docker.sock",
        "version": "v1.22"
    },
    "listen": ":8080"
}
```

- Docker host is not required. Default is `unix:///var/run/docker.sock`
- Docker version is not required. Default is empty string.

## Supported Services

- Bitbucket (should be pull request webhook)
- GitHub (should be a pull request webhook)
- GitLab (should be a merge request webhook)
- JIRA (should be a webhook when a issue is updated)

Only merged pull request will be handled. In JIRA you can set a scope for your webhook.

## Run in Docker

```
docker build -t farmer .
docker run -i -t -d -e VIRTUAL_HOST=test.example -v /var/run/docker.sock:/var/run/docker.sock:ro farmer
```

## License

MIT © Isotop