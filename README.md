# Logger

Simple Logger for Go applications. This is a simple logger, built on top of standar package `log` that added the option of logging
information with three log levels: `DEBUG`, `INFO` and `ERROR`.

It also has a small integration with Slack, where error can be automatically sent to a pre-configured slack channel, using webhooks.
