# Logging

Logging is a package for logging and encoding various [go-metrics-plus](https://github.com/zeim839/go-metrics-plus) metrics. The package may be used to log metrics to stdout through the use of an Encoder interface, which transforms metrics into plain text.

The package has built-in encoders for graphite plain text, prometheus expositional format, and Stasd line protocol.
