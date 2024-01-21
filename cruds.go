package main

import (
	"log/slog"
	"main/src"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func init() {
	// logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	src.Logger = slog.New(handler)
}

func main() {
	// config
	f, err := os.Open("config.yml")
	if err != nil {
		src.Logger.Error("Cant open config yaml", err)
		return
	}
	defer f.Close()

	var cfg src.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		src.Logger.Error("Cant decode config yaml", err)
		return
	}

	// router
	router := src.SetupRouter()

	// server
	src.Logger.Info("server " + cfg.Server.Port + " is started")
	err = http.ListenAndServe(":"+cfg.Server.Port, router)
	if err != nil {
		src.Logger.Error("Problems with server: ", err)
	}
}
