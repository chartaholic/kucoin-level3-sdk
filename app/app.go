package app

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/Kucoin/kucoin-level3-sdk/api"
	"github.com/Kucoin/kucoin-level3-sdk/builder"
	"github.com/Kucoin/kucoin-level3-sdk/events"
	"github.com/Kucoin/kucoin-level3-sdk/service"
	"github.com/Kucoin/kucoin-level3-sdk/utils/log"
)

type App struct {
	apiService *kucoin.ApiService
	symbol     string

	enableOrderBook bool
	level3Builder   *builder.Builder

	enableEventWatcher bool
	eventWatcher       *events.Watcher
	redisPool          *service.Redis

	rpcPort  string
	rpcToken string
}

func NewApp(symbol string, rpcPort string, rpcKey string) *App {
	if symbol == "" {
		panic("symbol is required")
	}

	if rpcPort == "" {
		panic("rpcPort is required")
	}

	if rpcKey == "" {
		panic("rpckey is required")
	}

	apiService := kucoin.NewApiServiceFromEnv()
	level3Builder := builder.NewBuilder(apiService, symbol)

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")
	var redisDBEnv = os.Getenv("REDIS_DB")
	var redisDB = 0
	if redisDBEnv != "" {
		redisDB, _ = strconv.Atoi(redisDBEnv)
	}
	redisPool := service.NewRedis(redisHost, redisPassword, redisDB, rpcKey, symbol, rpcPort)

	eventWatcher := events.NewWatcher(redisPool)

	return &App{
		apiService: apiService,
		symbol:     symbol,

		enableOrderBook: os.Getenv("ENABLE_ORDER_BOOK") == "true",
		level3Builder:   level3Builder,

		enableEventWatcher: os.Getenv("ENABLE_EVENT_WATCHER") == "true",
		redisPool:          redisPool,
		eventWatcher:       eventWatcher,

		rpcPort:  rpcPort,
		rpcToken: os.Getenv("RPC_TOKEN"),
	}
}

func (app *App) Run() {
	if app.enableOrderBook {
		go app.level3Builder.ReloadOrderBook()
	}

	if app.enableEventWatcher {
		go app.eventWatcher.Run()
	}

	//rpc server
	go api.InitRpcServer(app.rpcPort, app.rpcToken, app.level3Builder, app.eventWatcher)

	app.websocket()
}

func (app *App) writeMessage(msgRawData json.RawMessage) {
	//log.Info("raw message : %s", kucoin.ToJsonString(msgRawData))
	if app.enableOrderBook {
		app.level3Builder.Messages <- msgRawData
	}

	if app.enableEventWatcher {
		app.eventWatcher.Messages <- msgRawData
	}

	const msgLenLimit = 50
	if len(app.level3Builder.Messages) > msgLenLimit ||
		len(app.eventWatcher.Messages) > msgLenLimit {
		log.Error(
			"msgLenLimit: app.level3Builder.Messages: %d, app.eventWatcher.Messages: %d, app.verify.Messages: %d",
			len(app.level3Builder.Messages),
			len(app.eventWatcher.Messages),
		)
	}
}

func (app *App) websocket() {
	//todo recover dingTalk ?
	apiService := app.apiService

	rsp, err := apiService.WebSocketPublicToken()
	if err != nil {
		panic(err)
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		panic(err)
	}

	c := apiService.NewWebSocketClient(tk)

	mc, ec, err := c.Connect()
	if err != nil {
		panic(err)
	}

	ch := kucoin.NewSubscribeMessage("/market/level3:"+app.symbol, false)
	if err := c.Subscribe(ch); err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-ec:
			c.Stop()
			panic(err)

		case msg := <-mc:
			app.writeMessage(msg.RawData)
		}
	}
}
