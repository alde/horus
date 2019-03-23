package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/braintree/manners"
	"github.com/sirupsen/logrus"

	"github.com/alde/horus/backend/encryptor"

	"github.com/alde/horus/backend/config"
	"github.com/alde/horus/backend/database"
	"github.com/alde/horus/backend/server"
	"github.com/alde/horus/backend/version"
)

func main() {
	go catchInterrupt()

	configFile := flag.String("config", "", "Specify a config.toml file")
	flag.Parse()

	cfg := config.Initialize(*configFile)
	setupLogging(cfg)
	ctx := context.Background()

	db, err := database.Init(ctx, cfg)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create database")
	}

	bind := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	logrus.WithFields(logrus.Fields{
		"version": version.Version,
		"address": cfg.Server.Address,
		"port":    cfg.Server.Port,
	}).Info("launching Horus backend")
	router := server.NewRouter(cfg, db, encryptor.NewGoogleCloudKMS(ctx, cfg))
	if err := manners.ListenAndServe(bind, router); err != nil {
		logrus.WithError(err).Fatal("unrecoverable error")
	}
}

func setupLogging(cfg *config.Config) {
	if cfg.Logging.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	if s != os.Interrupt && s != os.Kill {
		return
	}
	logrus.Info("shutting down")
	os.Exit(0)
}
