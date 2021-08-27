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
docker run --rm --name some-postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=user -p 5432:5432 -d postgres
```

### Start a Redis instance
```
docker run --rm --name some-redis -p 6379:6379 -d redis
```

### Example configuration (REST+gRPC server)

```yaml
rest:
  port: 9988
grpc:
  port: 50051
#dm-cache:
  #type: olric
  #mode: lan
  type: redis
  host-address: 127.0.0.1:6379
  password: ""
db-cache:
  driver-name: postgres
  data-source-name: "host=localhost port=5432 user=user password=secret dbname=user sslmode=disable"
  local-cache-size: 0
  historicize:
    # example: run every 1 min - see https://crontab.guru/#*_*_*_*_*
    schedule: "* * * * *"
    #format: csv
    format: parquet
    tmp-dir: "./"
    date-partitioner: "year=2006/month=01/day=02"
    delete-local: true
    s3:
      endpoint: "play.min.io"
      use-ssl: false
      bucket: mytestbucket
      access-key-varname: ACCESSKEY
      secret-key-varname: SECRETKEY
```

with the date partitioner format `year=2006/month=01/day=02` referring to the Golang's year, month, day format as also described [here](https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format).

The `delete-local: bool` defines whether to remove local partitions upon file upload.

### Run the service

```bash
❯ ./igovium --config conf.yaml
```

### REST API example

PUT on `http://localhost:9988`:
```json
{
    "key":"mykey",
    "value": {"myvalue":1, "myotherval":100},
    "ttl" : "1h"
}
```

Returns 200 OK and the json payload.

GET on `http://localhost:9988/mykey`:
```json
{
    "myvalue": 1,
    "myotherval": 100
}
```

### Run the grpc client example
Please find an example gRPC client [here](examples/grpc_client/client.go).

```bash
❯ ./examples/grpc_client/grpc_client
2021/08/27 16:26:46 putting: k='key', v='{"mykey":"this-is-my-test-value"}'
2021/08/27 16:26:46 put response: res='', err='<nil>'
2021/08/27 16:26:46 get response: value:"{\"mykey\":\"this-is-my-test-value\"}"
```

### Historicizer example using play.minio S3
Here the `conf.yaml` settings:
```yaml
db-cache:
  historicize:
    s3:
      endpoint: "play.min.io"
      use-ssl: false
      bucket: mytestbucket
      access-key-varname: ACCESSKEY
      secret-key-varname: SECRETKEY
```

Here the actual variable containing the access key and secret key:
```bash
export ACCESSKEY=Q3AM3UQ867SPQQA43P2F
export SECRETKEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

