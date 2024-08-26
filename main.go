package main

import (
	"encoding/json"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/go-vgo/robotgo"
	log "go.uber.org/zap"
	"strings"
)

const version = "0.0.1-beta"

var logger *log.Logger

func main() {
	rawJSON := []byte(
		fmt.Sprintf(`{
			"level": "debug",
			"encoding": "json",
			"outputPaths": ["stdout", "/tmp/amfa_logs"],
			"errorOutputPaths": ["stderr"],
			"initialFields": {"amfa": "%s"},
			"encoderConfig": {
				"messageKey": "message",
				"levelKey": "level",
				"levelEncoder": "lowercase"
			}
		}`, version))

	var cfg log.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger = log.Must(cfg.Build())
	defer func(logger *log.Logger) {
		err := logger.Sync()
		if err != nil && !strings.Contains(
			err.Error(),
			"inappropriate ioctl for device",
		) {
			logger.Error("failed to sync logger", log.Error(err))
		}
	}(logger)

	logger.Info("logger construction succeeded")
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle(fmt.Sprintf("amfa %s", version))
	systray.SetTooltip("keeps me awake!!!")

	about := systray.AddMenuItem("About", "What's amfa?")
	go func() {
		<-about.ClickedCh
		systray.SetTooltip("Snort it.")
		robotgo.Alert(fmt.Sprintf(
			"amfa %s",
			version,
		),
			"amfa is a simple app that keeps your computer awake by faking mouse movements.",
			"kewl!",
		)
	}()

	mQuit := systray.AddMenuItem("Quit", "Let's get some sleep...")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	logger.Info("I'm out!")
}
