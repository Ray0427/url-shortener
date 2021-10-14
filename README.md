# url-shortener

> url shortener RESTful API

## Requiredment

- Go - greater than 1.13

- Mysql - for persistent data
- Redis - for cache data

## Environment Variable

Copy dotenv example

```shell
cp .env.example .env
```

Fulfill your config for Mysql, Redis and hashids

## Usage

### Download go module

```shell
go mod download
```

### Unit test

```shell
go test -v -cover=true ./...  
```

### Run server

```shell
go run main.go
```

## API

### Upload URL API

```shell
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{
"url": "<original_url>",
"expireAt": "2021-02-08T09:20:41Z"
}'
```

#### response

```json
{
  "id": "<url_id>",
  "shortUrl": "http://localhost/<url_id>"
}
```

### Redirect URL API

```shell
curl -L -X GET http://localhost/<url_id>
```

#### response

 REDIRECT to original URL

## TO DO

- [ ] unit test
- [ ] validator
- [ ] flow diagram