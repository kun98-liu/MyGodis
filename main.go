/*
main.go文件是该项目的入口。
实现功能如下：
1.读取配置文件
2.启动Server
*/

package main

import (
	"fmt"
	"os"

	"github.com/kun98-liu/MyGodis/config"
	"github.com/kun98-liu/MyGodis/lib/logger"
	RedisServer "github.com/kun98-liu/MyGodis/redis/server"
	"github.com/kun98-liu/MyGodis/tcp"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

var defaultProperties = &config.ServerProperties{}

func main() {

	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "Godis",
		Ext:        "log",
		TimeFormat: "2006-01-01",
	})

	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {

		if fileExist("redis.conf") {
			config.SetupConfig("redis.config")
		} else {
			config.Properties = defaultProperties
		}

	} else {
		config.SetupConfig(configFilename)
	}

	tcpConfig := &tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}
	err := tcp.ListenAndServeWithSignal(tcpConfig, RedisServer.MakeHandler())

	if err != nil {
		logger.Error(err)
	}

}

func fileExist(s string) bool {
	info, err := os.Stat(s)
	return err == nil && !info.IsDir()
}
