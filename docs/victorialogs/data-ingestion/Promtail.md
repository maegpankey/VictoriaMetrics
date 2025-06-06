---
weight: 4
title: Promtail setup
disableToc: true
menu:
  docs:
    parent: "victorialogs-data-ingestion"
    weight: 4
tags:
  - logs
aliases:
  - /victorialogs/data-ingestion/Promtail.html
  - /victorialogs/data-ingestion/promtail.html
---
[Promtail](https://grafana.com/docs/loki/latest/clients/promtail/), [Grafana Agent](https://grafana.com/docs/agent/latest/)
and [Grafana Alloy](https://grafana.com/docs/alloy/latest/) are default log collectors for Grafana Loki.
They can be configured to send the collected logs to VictoriaLogs according to the following docs.

Specify [`clients`](https://grafana.com/docs/loki/latest/clients/promtail/configuration/#clients) section in the configuration file
for sending the collected logs to [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/):

```yaml
clients:
  - url: "http://localhost:9428/insert/loki/api/v1/push"
```

Substitute `localhost:9428` address inside `clients` with the real TCP address of VictoriaLogs.

VictoriaLogs automatically parses JSON string from the log message into [distinct log fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model).
This behavior can be disabled by passing `-loki.disableMessageParsing` command-line flag to VictoriaLogs or by adding `disable_message_parsing=1` query arg
to the `/insert/loki/api/v1/push` url in the config of log shipper:

```yaml
clients:
  - url: "http://localhost:9428/insert/loki/api/v1/push?disable_message_parsing=1"
```

In this case the JSON with log fields is stored as a string in the [`_msg` field](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field),
so later it could be parsed at query time with the [`unpack_json` pipe](https://docs.victoriametrics.com/victorialogs/logsql/#unpack_json-pipe).
JSON parsing at query can be slow and can consume a lot of additional CPU time and disk read IO bandwidth. That's why it is
recommended enabling JSON message parsing at data ingestion.

VictoriaLogs uses [log stream labels](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields) defined at the client side,
e.g. at Promtail, Grafana Agent or Grafana Alloy. Sometimes it may be needed overriding the set of these fields. This can be done via `_stream_fields`
query arg. For example, the following config instructs using only the `instance` and `job` labels as log stream fields, while other labels
will be stored as [usual log fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model):

```yaml
clients:
  - url: "http://localhost:9428/insert/loki/api/v1/push?_stream_fields=instance,job"
```

It is recommended verifying whether the initial setup generates the needed [log fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model)
and uses the correct [stream fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields).
This can be done by specifying `debug` [parameter](https://docs.victoriametrics.com/victorialogs/data-ingestion/#http-parameters)
and inspecting VictoriaLogs logs then:

```yaml
clients:
  - url: "http://localhost:9428/insert/loki/api/v1/push?debug=1"
```

If some [log fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) must be skipped
during data ingestion, then they can be put into `ignore_fields` [parameter](https://docs.victoriametrics.com/victorialogs/data-ingestion/#http-parameters).
For example, the following config instructs VictoriaLogs to ignore `filename` and `stream` fields in the ingested logs:

```yaml
clients:
  - url: 'http://localhost:9428/insert/loki/api/v1/push?ignore_fields=filename,stream'
```

See also [these docs](https://docs.victoriametrics.com/victorialogs/data-ingestion/#http-parameters) for details on other supported query args.
There is no need in specifying `_time_field` query arg, since VictoriaLogs automatically extracts timestamp from the ingested Loki data.

By default the ingested logs are stored in the `(AccountID=0, ProjectID=0)` [tenant](https://docs.victoriametrics.com/victorialogs/#multitenancy).
If you need storing logs in other tenant, then specify the needed tenant via `tenant_id` field
in the [Loki client configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/#clients)
The `tenant_id` must have `AccountID:ProjectID` format, where `AccountID` and `ProjectID` are arbitrary uint32 numbers.
For example, the following config instructs VictoriaLogs to store logs in the `(AccountID=12, ProjectID=34)` [tenant](https://docs.victoriametrics.com/victorialogs/#multitenancy):

```yaml
clients:
  - url: "http://localhost:9428/insert/loki/api/v1/push"
    tenant_id: "12:34"
```

The ingested log entries can be queried according to [these docs](https://docs.victoriametrics.com/victorialogs/querying/).

See also [data ingestion troubleshooting](https://docs.victoriametrics.com/victorialogs/data-ingestion/#troubleshooting) docs.
