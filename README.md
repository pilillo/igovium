# IgoVIUM

Multi-level caching service in Go.

## Example

### Start a Postgres instance

```
docker run --rm --name some-postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=user -p 5432:5432 postgres
```

```yaml
port: 9988
dm-cache:
  mode: lan
db-cache:
  driver-name: postgres
  data-source-name: "host=localhost port=5432 user=user password=secret dbname=user sslmode=disable"
  local-cache-size: 0
```