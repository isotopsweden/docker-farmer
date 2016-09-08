# docker-farmer

Simple Go project that will handle payloads from different services and remove Docker containers based on the payload and the configured domain.

## Config

Example configuration:

```
{
    "domain": "test.example.com",
    "docker": {
        "host": "unix:///var/run/docker.sock",
        "version": "v1.22"
    },
    "listen": ":8080"
}
```

## Supported Services

- Bitbucket (should be pull request webhook)
- GitHub (should be a pull request webhook)
- GitLab (should be a merge request webhook)
- JIRA (should be a webhook when a issue is updated)

Only merged pull request will be handled. In JIRA you can set a scope for your webhook.

## License

MIT © Isotop