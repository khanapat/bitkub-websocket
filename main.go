package main

import (
	"bitkub-websocket/bitkub"
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	streams := []string{
		bitkub.THBBTCTicker,
		// bitkub.THBETHTicker,
	}

	socketUrl := fmt.Sprintf("%s/%s", viper.GetString("client.bitkub.websocket"), strings.Join(streams, ","))
	c, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		logger.Error(err.Error())
	}
	defer c.Close()

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
			var tickerMarketWebsocket bitkub.TickerMarketWebsocket
			if err := json.Unmarshal(message, &tickerMarketWebsocket); err != nil {
				logger.Error(err.Error())
				continue
			}
			switch tickerMarketWebsocket.Stream {
			case bitkub.THBBTCTicker:
				logger.Info(fmt.Sprintf("---- %s ----", tickerMarketWebsocket.Stream))
				logger.Info(fmt.Sprintf("Last: %f", tickerMarketWebsocket.Last))
				logger.Info(fmt.Sprintf("Payload: %s", message))
			case bitkub.THBETHTicker:
				logger.Info(fmt.Sprintf("Payload: %s", message))
			default:
			}
		}
	}
}

func initViper() {
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

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
