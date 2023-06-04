# ksqlite

Kafka connector for reading/writing from/to SQLite

### Examples

- Reader
- Writer

### Testing locally

#### Install docker
TODO

#### Deploy infrastructure dependencies

Run docker-compose command
```bash
docker-compose -f deployments/docker-compose.yaml up -d
```

#### Run integration tests

```bash
go test -v ./...
```

#### Stop infrastructure services

```bash
docker-compose -f deployments/docker-compose.yaml down 
```