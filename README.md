# url-shortener

> url shortener RESTful API

## Description

這個服務主要會使用 Mysql 作為 database、Redis作為cache storage，在DB的選擇上，這個情境並沒有特別的需求，所以選擇目前最熟悉的Mysql，為了降低DB讀取量、以及提高response的速度，所以加上了redis 作為server cache

### url_id的設計

短連結的服務通常會避免連結很容易被試出來，所以將url_id會使用random string的方式，但是在Mysql使用random string作為Primary key會遇到一些問題，因為Mysql的Clustered Index 是採用 pk 建立的 B+Tree，不連續的row insert會導致InnoDB 必須要計算出適當的位置來安插新的資料，paging的資料需要重新查找或是移出快取，所以這邊使用auto increment的pk，搭配上 `go-hashids`的lib，將integer 的 PK轉為hash string

### create url流程

1. 當server收到request，首先會做格式的檢查
2. 在DB寫入這筆url與expireAt
3. 將這筆row的pk轉回hashID
4. 以hashID作為key，url與expireAt轉為JSONString後作為value，存入redis，ttl設為1小時
5. 回傳使用者hashId與shortUrl

### get url流程

1. 當server收到request，首先會去redis尋找是否有這筆資料的cache？
   1. 如果有資料但expireAt已過，回傳Not Found
   2. 如果有資料且未過期，redirect to url
   3. 如果有資料但內容是null字串，回傳Not Found
   4. 如果沒有資料，進到下一步
2. 將hashID轉回row的pk
3. 在DB查詢是否有這筆資料？
   1. 有資料
      1. 以hashID作為key，url與expireAt轉為JSONString後作為value，存入redis，ttl設為1小時
      2. 是否過期？
         1. 未過期，redirect to url
         2. xpireAt已過，回傳Not Found
   2. 無資料
      1. 以hashID作為key，將null字串作為value，存入redis，ttl設為1小時
      2. 回傳Not Found

### 其他lib

- github.com/joho/godotenv - 將`.env`讀取放入 environment variable

- github.com/caarlos0/env/v6 - 將 environment variable 塞入config struct
- github.com/go-redis/redis/v8 - redis driver
- gorm.io/gorm - orm for database
- github.com/speps/go-hashids/v2 - encode & decode between integer and hashId

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
HASHID_SALT為字串, HASHID_MIN_LENGTH為url_id長度

## Usage

### Build

```shell
go build 
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
- [x] cache unit test
- [x] dao unit test
- [x] controller unit test
