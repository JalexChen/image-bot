# image-bot

The goal is to standardize and automate image manifests for convenience and vm images. These manifest values can be used with automation software like Ansible to provision an instance

# Getting Started

This bot runs shell commmands, specifically for ubuntu/debian, so you'll need to run this on one of those containers.
You'll also need an environment variable called `GITHUB_TOKEN` as this bot utilizes the Github API. You can include this in an `.env` file

Use, add, or modify the files in `testdata` for the configuration and manifest

Run `docker compose up -d` then run:

```
go run main.go -c << config filepath >> -m << manifest filepath >>

### example
go run main.go -c testdata/jammy-config.yml -m testdata/jammy-manifest.yml

```

# Future

- Support more shells/os
- Support more undefined sources
- Support sdkmanager (for Android)
