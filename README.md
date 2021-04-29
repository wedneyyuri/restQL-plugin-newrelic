# NewRelic Plugin for RestQL

This plugin is an interface that allow [RestQL](https://github.com/b2wdigital/restQL-golang) to use [New Relic](https://newrelic.com/) as Application Performance Monitoring solution.

## Building

To use this plugin you need the [RestQL CLI]() tool. You can build a custom RestQL binary with this plugin using the following command:

```shell
$ RUN restQL-cli build --with github.com/wedneyyuri/restQL-plugin-newrelic@v6.0.0 --output ./restQL v6.2.0
```

## Usage

This plugin can be configured [through environment variables](https://github.com/newrelic/go-agent/blob/198c033a21ef66200032d444e1be06b354175516/v3/newrelic/config_options.go#L54).

### Minimal Configuration:

- `NEW_RELIC_APP_NAME`: (**REQUIRED**) sets AppName
- `NEW_RELIC_LICENSE_KEY`: (**REQUIRED**) sets License

### Other configurations

- `NEW_RELIC_ATTRIBUTES_EXCLUDE`: sets Attributes.Exclude using a comma-separated list, eg. "request.- `headers.host,request.method"
- `NEW_RELIC_ATTRIBUTES_INCLUDE`: sets Attributes.Include using a comma-separated list
- `NEW_RELIC_DISTRIBUTED_TRACING_ENABLED`: sets DistributedTracer.Enabled using [strconv.ParseBool](https://golang.org/pkg/strconv/#ParseBool)
- `NEW_RELIC_ENABLED`: sets Enabled using [strconv.ParseBool](https://golang.org/pkg/strconv/#ParseBool)
- `NEW_RELIC_HIGH_SECURITY`: sets HighSecurity using [strconv.ParseBool](https://golang.org/pkg/strconv/#ParseBool)
- `NEW_RELIC_HOST`: sets Host
- `NEW_RELIC_INFINITE_TRACING_SPAN_EVENTS_QUEUE_SIZE`: sets InfiniteTracing.SpanEvents.QueueSize using [strconv.Atoi](https://golang.org/pkg/strconv/#Atoi).
- `NEW_RELIC_INFINITE_TRACING_TRACE_OBSERVER_PORT`: sets InfiniteTracing.TraceObserver.Port using [strconv.Atoi](https://golang.org/pkg/strconv/#Atoi).
- `NEW_RELIC_INFINITE_TRACING_TRACE_OBSERVER_HOST`: sets InfiniteTracing.TraceObserver.Host
- `NEW_RELIC_LABELS`: sets Labels using a semi-colon delimited string of colon-separated - `pairs, eg. "Server:One;DataCenter:Primary"
- `NEW_RELIC_LOG`: sets Logger to log to either "stdout" or "stderr" (filenames are not - `supported)
- `NEW_RELIC_LOG_LEVEL`: controls the NEW_RELIC_LOG level, must be "debug" for debug, or - `empty for info
- `NEW_RELIC_PROCESS_HOST_DISPLAY_NAME`: sets HostDisplayName
- `NEW_RELIC_SECURITY_POLICIES_TOKEN`: sets SecurityPoliciesToken
- `NEW_RELIC_UTILIZATION_BILLING_HOSTNAME`: sets Utilization.BillingHostname
- `NEW_RELIC_UTILIZATION_LOGICAL_PROCESSORS`: sets Utilization.LogicalProcessors using [strconv.Atoi](https://golang.org/pkg/strconv/#Atoi).
- `NEW_RELIC_UTILIZATION_TOTAL_RAM_MIB`: sets Utilization.TotalRAMMIB using [strconv.Atoi](https://golang.org/pkg/strconv/#Atoi).

## License

The [MIT license](https://mit-license.org/). See the LICENSE file.
