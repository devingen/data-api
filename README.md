# Veri API

Contains Veri API server.

## Deployment

### Docker

Docker build needs to access private repositories under github.com/devingen.

To let it pull the repos, you need to [generate a personal access token](https://github.com/settings/tokens)
and pass GIT_TOKEN to the command as follows.

```shell
GIT_TOKEN=GITHUB_TOKEN_GENERATED_ON_WEBSITE
make release-docker GIT_TOKEN=$GIT_TOKEN IMAGE_TAG=0.1.9
make release-docker GIT_TOKEN=$GIT_TOKEN IMAGE_TAG=latest
```