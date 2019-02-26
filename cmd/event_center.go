package cmd

import (
	"fmt"
	"os"

	"github.com/sipt/shuttle/config"
	. "github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/log"
)

var eventChan chan *EventObj

func ListenEvent() {
	eventChan = make(chan *EventObj, 1)
	go dealEvent(eventChan)
}

func dealEvent(c chan *EventObj) {
	for {
		t := <-c
		switch t.Type {
		case EventShutdown.Type:
			log.Logger.Info("[Shuttle] is Shutdown, see you later!")
			Shutdown(config.CurrentConfig().General.SetAsSystemProxy)
			os.Exit(0)
			return
		case EventReloadConfig.Type:
			_, err := reloadConfig(config.CurrentConfigFile(), StopSocksSignal, StopHTTPSignal)
			if err != nil {
				log.Logger.Error("Reload Config failed: ", err)
				fmt.Println(err.Error())
				os.Exit(1)
			}
		case EventRestartHttpProxy.Type:
			StopHTTPSignal <- true
			go HandleHTTP(config.CurrentConfig(), StopHTTPSignal)
		case EventRestartSocksProxy.Type:
			StopSocksSignal <- true
			go HandleSocks5(config.CurrentConfig(), StopSocksSignal)
		case EventRestartController.Type:
			controller.ShutdownController()
			go controller.StartController(config.CurrentConfig(), eventChan)
		}
	}
}

func ReloadConfig(configPath string) int {
	_, err := reloadConfig(configPath, StopSocksSignal, StopHTTPSignal)
	if err != nil {
		return 1
	}
	return 0
}
