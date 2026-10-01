package main

import (
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/joho/godotenv"

	"micro-blog/internal/app"
	"micro-blog/internal/logger"
	"micro-blog/internal/validator"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	application := app.InitializeApp()
	defer application.DB.Close()

	if application.Env == logger.Dev {
		go pprof(application)
	}

	port := string(application.Port)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           application.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	application.Logger.Info("server started", "port", port)

	if err := server.ListenAndServe(); err != nil {
		application.Logger.Error("server failed", "error", err)
		return err
	}

	return nil
}

func pprof(application *app.App) {
	debugServer := &http.Server{
		Addr:              "localhost:6060",
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	application.Logger.Info("pprof started", "addr", "localhost:6060")

	if err := debugServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		application.Logger.Error("pprof server failed", "error", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	initValidator()
}

func initValidator() {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	validator.Translator, _ = uni.GetTranslator("en")

	err := entranslations.RegisterDefaultTranslations(validator.Validate, validator.Translator)
	if err != nil {
		return
	}
}
