package main

import (
	"bitkub-websocket/bitkub"
	"bitkub-websocket/common"
	"bitkub-websocket/internal/redis"
	"bitkub-websocket/logz"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initViper()
	initTimezone()
}

func main() {
	logger := logz.NewLogConfig()

	pool := redis.NewRedisConn()
	defer pool.Close()

	setDataNoExpireRedisFn := redis.NewSetDataNoExpireRedisFn(pool)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	streams := []string{
		bitkub.THBBTCTicker,
		bitkub.THBETHTicker,
		// bitkub.THBBTCTrade,
	}

	socketUrl := fmt.Sprintf("%s/%s", viper.GetString("client.bitkub.websocket"), strings.Join(streams, ","))
	c, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		logger.Error(err.Error())
	}
	defer c.Close()

	logger.Info(fmt.Sprintf("Websocket Starting with stream channel: %s", strings.Join(streams, ",")))

	for {
		select {
		case <-signals:
			logger.Info("Stopping...")
			return
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Error(err.Error())
				return
			}
			var bitkubWebsocket bitkub.BitkubWebsocket
			if err := json.Unmarshal(message, &bitkubWebsocket); err != nil {
				logger.Error(err.Error())
				continue
			}

			switch {
			case strings.HasPrefix(bitkubWebsocket.Stream, bitkub.Ticker):
				var tickerMarketWebsocket bitkub.TickerMarketWebsocket
				if err := json.Unmarshal(message, &tickerMarketWebsocket); err != nil {
					logger.Error(err.Error())
				}
				logger.Info(fmt.Sprintf("---- %s ----", tickerMarketWebsocket.Stream))
				logger.Info(fmt.Sprintf("Last: %f", tickerMarketWebsocket.Last))
				logger.Info(fmt.Sprintf("Open: %f | Close: %f", tickerMarketWebsocket.Open, tickerMarketWebsocket.Close))
				if err := setDataNoExpireRedisFn(common.MapRateTokenRedis[tickerMarketWebsocket.Stream], tickerMarketWebsocket.Last); err != nil {
					logger.Error(err.Error())
				}
			case strings.HasPrefix(bitkubWebsocket.Stream, bitkub.Trade):
				var tradeMarketWebsocket bitkub.TradeMarketWebsocket
				if err := json.Unmarshal(message, &tradeMarketWebsocket); err != nil {
					logger.Error(err.Error())
				}
				logger.Info(fmt.Sprintf("---- %s ----", tradeMarketWebsocket.Stream))
				logger.Info(fmt.Sprintf("Rate: %f with Amount: %f", tradeMarketWebsocket.Rate, tradeMarketWebsocket.Amount))
			default:
			}
			logger.Info(fmt.Sprintf("Payload: %s", string(message)))
		}
	}
}

func initViper() {
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

	viper.SetDefault("redis.max-idle", 3)
	viper.SetDefault("redis.timeout", "60s")
	// viper.SetDefault("redis.host", "docker.for.mac.localhost:6379")
	viper.SetDefault("redis.host", "localhost:6379")
	viper.SetDefault("redis.password", "P@ssw0rd")

	viper.SetDefault("client.bitkub.websocket", "wss://api.bitkub.com/websocket-api")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("error loading location 'Asia/Bangkok': %v\n", err)
	}
	time.Local = ict
}
