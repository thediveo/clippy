/*
Package debug supplies the “--debug” and “--tinted” CLI flags for configuring
the default structured logger.

  - “--debug” enables structured logging to stderr from debug level on and
    upwards (so no tracing).
  - “--tinted” enables tinted logs using the [lmittmann/tint] module.

[lmittmann/tint]: https://github.com/lmittmann/tint
*/
package debug
