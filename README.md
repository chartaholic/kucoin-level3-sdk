# Kucoin Level3 market

## guide
  [中文文档](README_CN.md)

## Installation

1. install dependencies

```
go get github.com/JetBlink/orderbook
go get github.com/go-redis/redis
go get github.com/gorilla/websocket
go get github.com/joho/godotenv
go get github.com/Kucoin/kucoin-go-sdk
go get github.com/shopspring/decimal
```

2. build

```
CGO_ENABLED=0 go build -ldflags '-s -w' -o kucoin_market kucoin_market.go
```

or you can download the latest available [release](https://github.com/Kucoin/kucoin-level3-sdk/releases)

## Usage

1. [vim .env](.env):
    ```
    # API_SKIP_VERIFY_TLS=1
    
    API_BASE_URI=https://api.kucoin.com
    
    # If open order book true otherwise false
    ENABLE_ORDER_BOOK=true
    
    # If open event watcher true otherwise false
    ENABLE_EVENT_WATCHER=true
    
    # Password for RPS calls. Pass the same when calling
    RPC_TOKEN=market-token
    
    REDIS_HOST=127.0.0.1:6379
    REDIS_PASSWORD=
    REDIS_DB=
    ```

1. Run Command：

    ```
    ./kucoin_market -c .env -symbol BTC-USDT -p 9090 -rpckey BTC-USDT
    ```
    

## Docker Usage

1. Build docker image

   ```
   docker build -t kucoin_market .
   ```

1. [vim .env](.env)

1. Run

  ```
  docker run --rm -it -v $(pwd)/.env:/app/.env --net=host kucoin_market
  ```

## RPC Method

> endpoint : 127.0.0.1:9090
> the sdk rpc is based on golang jsonrpc 1.0 over tcp.

see:[python jsonrpc client demo](./demo/python-demo/level3/rpc.py)

* Get Part Order Book
    ```
    {"method": "Server.GetPartOrderBook", "params": [{"token": "your-rpc-token", "number": 1}], "id": 0}
    ```
    
* Get Full Order Book
    ```
    {"method": "Server.GetOrderBook", "params": [{"token": "your-rpc-token"}], "id": 0}
    ```

* Add Event ClientOids To Channels
    ```
    {"method": "Server.AddEventClientOidsToChannels", "params": [{"token": "your-rpc-token", "data": {"clientOid": ["channel-1", "channel-2"]}}], "id": 0}
    ```

* Add Event OrderIds To Channels
    ```
    {"method": "Server.AddEventOrderIdsToChannels", "params": [{"token": "your-rpc-token", "data": {"orderId": ["channel-1", "channel-2"]}}], "id": 0}
    ```
## Python-Demo

> the demo including orderbook display

see:[python use_level3 demo](./demo/python-demo/order_book_demo.py)
- Run order_book.py
    ```
    command: python order_book.py
    describe: display orderbook
    ```
