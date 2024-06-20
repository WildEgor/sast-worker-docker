SAST wrapper for Dockerfiles

## Requirements
- [Git](http://git-scm.com/)
- [Go >= 1.22](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Task](https://taskfile.dev/)
- [Air](https://github.com/air-verse/air)
- [Protoc]()

## Links
- https://github.com/aquasecurity/trivy
- https://github.com/hadolint/hadolint
- https://github.com/air-verse/air

## Project structure:
- cmd - Contains main.go as entrypoint
- dist - builds
- internal
  - adapters - adapters for linters, rpc and vulnerabilities checks
  - configs - configuration (ports, api path)
  - services - business logic
  - utils - shared utils
- scripts - contains sh scripts for linters

## How to run

### Generate proto (details in Taskfile.yml)
```shell
task proto
```

### Test commands
>Test command for hadolint:
>```shell
>docker run --rm -i -e HADOLINT_FORMAT=json hadolint/hadolint < Dockerfile
>```
>
>Test command for trivy
>```shell
>docker run -v /var/run/docker.sock:/var/run/docker.sock -v $HOME/Library/Caches:/root/.cache/ aquasec/trivy:0.52.2 image python:3.4-alpine
>```

### Before start
> First start trivy server!
>```shell
>docker-compose up -d --build trivy-server
>```

### Using docker

Run in dev mode
```shell
docker-compose up -d --build sast-worker-docker
```

Run in debug mode
```shell
docker-compose up -d --build sast-worker-docker-debug
```

### Using air

Run app with hot-reload (details in Taskfile.yml)
```shell
task local-dev
```
