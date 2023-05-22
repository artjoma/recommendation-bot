package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"os/signal"
	"recommendation-bot/services"
	"syscall"
	"time"
)

const (
	AppCfgLocation       = "config.yaml"
	TriggersCfgLocation  = "triggers.json"
	TemplatesCfgLocation = "templates.json"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	// init logging sys
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05", NoColor: true})
}

func main() {
	sigs := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Start read application config")
	// read config
	cData, err := os.ReadFile(AppCfgLocation)
	panicOnErr(err, "Failed read "+AppCfgLocation)
	cfg := &services.AppCfg{}
	err = yaml.Unmarshal(cData, cfg)
	panicOnErr(err, "Invalid json format of config.json")
	log.Info().Msg("End read application config")
	/*
		Init services!
	*/
	triggerService, err := services.NewTriggerService(TriggersCfgLocation)
	panicOnErr(err, "Failed to create Trigger service")

	templateService, err := services.NewTemplateService(TemplatesCfgLocation)
	panicOnErr(err, "Failed to create Template service")

	// init application context
	appCtx := services.NewAppCtx(cfg, triggerService, templateService)
	// init Telegram service
	telegramBotService := services.NewTelegramBot(appCtx)
	go telegramBotService.StartBot()
	// set active state to application
	appCtx.SetActiveState()

	// wait until OS process SIG shutdown
	<-sigs
	log.Info().Msg("Shutdown")
	appCtx.Destroy()
	log.Info().Msg("Exit from main()")
}
