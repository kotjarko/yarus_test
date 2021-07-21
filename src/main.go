package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const AppName = "YaRusTest"

func main() {
	// load conf
	cfg, err := parseConfig(AppName)
	if err != nil {
		logrus.Fatalf("unable to get config: %v",err)
	}

	// load currencies data
	base, err := InitCurrencyBase(cfg.Data.Url, cfg.Data.Retries, cfg.Data.RetryTimeout, cfg.Data.Timeout)
	if err != nil {
		logrus.Fatalf("unable to load currencies data: %v",err)
	}

	app := Application{
		Config:       cfg,
		CurrencyBase: base,
	}

	// run server
	callbackServer := http.Server{
		Addr:           app.Config.Web.Port,
		ReadTimeout:    app.Config.Web.ReadTimeout,
		WriteTimeout:   app.Config.Web.WriteTimeout,
		Handler:        app.createHandler(),
	}
	callbackErrors := make(chan error, 1)

	go func() {
		logrus.Infof("api started on %s", callbackServer.Addr)
		callbackErrors <- callbackServer.ListenAndServe()
	}()

	// check term signals and api errors
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-callbackErrors:
		logrus.Errorf("api error: %v", err)
	case <-osSignals:
		logrus.Info("starting shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), app.Config.Web.ShutdownTimeout)
		defer cancel()

		if err = callbackServer.Shutdown(ctx); err != nil {
			logrus.Errorf("graceful shutdown failed: %v", err)
			callbackServer.Close()
		}
	}
	logrus.Info("shutdown completed")
}
