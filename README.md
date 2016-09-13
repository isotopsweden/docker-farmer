# docker-farmer

Custom tool to remove docker containers and databases from services based on domain names.

## Config

Example configuration:

```js
{
    "domain": "test.example.com",
    "docker": {
        "host": "unix:///var/run/docker.sock",
        "version": "v1.22"
    },
    "database": {
        "prefix": "ab_",
        "container": "name or id",
        "type": "mysql"
    },
    "listen": ":80",
    "containers": {
        "exclude": ["redis"]
    }
}
```

- Docker host is not required. Default is `unix:///var/run/docker.sock`
- Docker version is not required. Default is empty string.

## Supported Services

- Bitbucket (should be pull request webhook)
- GitHub (should be a pull request webhook)
- GitLab (should be a merge request webhook or push events webhook)
- JIRA (should be a webhook when a issue is updated)

## Database

The database name should be a md5 hash of the domain for the container with the database prefix.

```
md5 <<<"abc-123.test.example.com"
5d4f0978db52bd0e588bef9cef98715f
```

Gives us

```
ab_5d4f0978db52bd0e588bef9cef98715f
```

## Run in Docker

```
docker build -t farmer .
docker run -i -t -d -e VIRTUAL_HOST=test.example.com -v /var/run/docker.sock:/var/run/docker.sock:ro farmer
```

## License

MIT © Isotop