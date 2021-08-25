# IgoVIUM

Multi-level caching service in Go.

Specifically: 
* [Distributed in-memory cache (L1)](cache/dm_cache_mapper.go)
* [DB-based cache (L2)](cache/db_cache.go)
* [Long term historization on persistent volumes (L3)](cache/historicizer.go)
  
Uses the following libraries:
* L1 - distributed in-memory cache
  * [Olric](https://github.com/buraksezer/olric)
  * [Redis](https://github.com/go-redis/redis)
* L2 - DB-based cache
  * [XORM](https://gitea.com/xorm/xorm) as ORM in Go targeting multiple DBs
* L3 - Historization to local and remote path

Historicizes to external volumes with any of the following formats:
* [CSV](cache/csv_formatter.go)
* [Parquet](cache/parquet_formatter.go)


## Example

### Start a Postgres instance

```
docker run --rm --name some-postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=user -p 5432:5432 postgres
```

### Example configuration (REST+gRPC server)

```yaml
rest:
  port: 9988
grpc:
  port: 50051
dm-cache:
  type: olric
  mode: lan
db-cache:
  driver-name: postgres
  data-source-name: "host=localhost port=5432 user=user password=secret dbname=user sslmode=disable"
  local-cache-size: 0
  historicize:
    # example: run every 1 min - see https://crontab.guru/#*_*_*_*_*
    schedule: "* * * * *"
    endpoint: "play.min.io"
    use-ssl: false
    bucket: mytestbucket
    format: csv
    partitioner: ""
    tmp-dir: "./"
```

### Run the service

```bash
❯ ./igovium --config conf.yaml
```

### REST API example

PUT on `http://localhost:9988`:
```json
{
    "key":"mykey",
    "value": {"myvalue":1},
    "ttl" : "1h"
}
```

Returns 200 OK and the json payload.

GET on `http://localhost:9988/mykey`:
```json
{
    "myvalue": 1
}
```

### Run the grpc client example
Please find an example gRPC client [here](examples/grpc_client/client.go).

```bash
❯ ./examples/grpc_client/grpc_client
2021/08/11 16:57:55 put response: 
2021/08/11 16:57:55 get response: value:"\x08\n\x00\x05value"
```
