package restqlnewrelic

import (
	"fmt"
	"os"
	"strconv"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const transactionEventsMaxSamples = "NEW_RELIC_TRANSACTION_EVENTS_MAX_SAMPLES_STORED"

func ExtraConfigFromEnvironment() newrelic.ConfigOption {
	return extraConfigFromEnvironment(os.Getenv)
}

func extraConfigFromEnvironment(getenv func(string) string) newrelic.ConfigOption {
	return func(cfg *newrelic.Config) {
		// Because fields could have been assigned in a previous
		// ConfigOption, we only want to assign fields using environment
		// variables that have been populated.  This is especially
		// relevant for the string case where no processing occurs.
		assignInt := func(field *int, name string) {
			if env := getenv(name); env != "" {
				if i, err := strconv.Atoi(env); nil != err {
					cfg.Error = fmt.Errorf("invalid %s value: %s", name, env)
				} else {
					*field = i
				}
			}
		}

		assignInt(&cfg.TransactionEvents.MaxSamplesStored, transactionEventsMaxSamples)
	}
}
