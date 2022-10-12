## ZAP Production logger

This package includes a Zap logger with an optional Loki backend. It operates in two modes:

- DEV
- PROD

In the `DEV` mode, only a `console` logger is used. Loki won't receive any notifications.

In the `PROD` mode, each message is written BOTH to console and Loki in a JSON mode

### Using Loki

Loki implements a notion of labels to build a hierarchical aggregation for logs. It is **required to keep the number of labels to the ABSOLUTE MINIMUM** to work effectively with Loki.

### Use cases

1. Use two labels: `env` and `service`. This will allow to filter the log as per TWO labels.
   Example:

```go
logger.Infow("First entry", "env", "prod", "service", "front")
```

This adds two labels and Loki will display this message for two filters

2. Use a single label

```go
logger.Infow("Second entry", "env", "prod")
```

This adds only one label - `env`, so Loki will display this entry only for `env = prod` setup
