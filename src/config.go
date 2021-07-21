package main

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type config struct {
	Debug	bool	`default:"false" envconfig:"YARUSTEST_DEBUG"`

	Data struct{
		Url     string        		`default:"https://www.cbr-xml-daily.ru/daily_json.js" envconfig:"YARUSTEST_DATA_URL"`
		Timeout time.Duration 		`default:"180s" envconfig:"YARUSTEST_DATA_TIMEOUT"`
		RetryTimeout time.Duration 	`default:"10s" envconfig:"YARUSTEST_DATA_RETRY_TIMEOUT"`
		Retries int           		`default:"3" envconfig:"YARUSTEST_DATA_RETRIES"`
	}

	Web struct {
		Port            string        `default:":80" envconfig:"YARUSTEST_WEB_PORT"`
		ReadTimeout     time.Duration `default:"120s" envconfig:"YARUSTEST_WEB_READ_TIMEOUT"`
		WriteTimeout    time.Duration `default:"120s" envconfig:"YARUSTEST_WEB_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `default:"5s" envconfig:"YARUSTEST_WEB_SHUTDOWN_TIMEOUT"`
	}
}

// load config by prefixes, not necessary here cuz constants filled in envconfig
// {app}_{config_category}_{param_name}
// example: YARUSTEST_DATA_TIMEOUT -> config.Data.Timeout
func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		_ = envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
