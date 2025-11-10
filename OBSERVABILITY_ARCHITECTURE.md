# Observability Architecture

## Overview

This Go TODO API implements a complete **LGTM observability stack** (Loki, Grafana, Tempo, MongoDB) with full trace-to-log correlation. The system provides:

- **Structured JSON logging** with trace context
- **Distributed tracing** with OpenTelemetry
- **Real-time dashboards** with metrics visualization
- **Bi-directional correlation** between traces and logs

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         User / Client                            │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP Requests
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Go TODO API (Port 8080)                      │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Middleware Chain:                                         │  │
│  │  • Logging Middleware    → Captures request/response      │  │
│  │  • Tracing Middleware    → Creates spans for requests     │  │
│  │  • CORS Middleware       → Handles cross-origin requests  │  │
│  │  • Auth Middleware       → API key validation             │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                   │
│  ┌─────────────────┐          ┌──────────────────────────────┐  │
│  │ Structured      │          │ OpenTelemetry SDK            │  │
│  │ Logger (slog)   │          │ • Creates traces & spans     │  │
│  │ • JSON format   │          │ • Propagates context         │  │
│  │ • trace_id      │          │ • Exports to Tempo via OTLP  │  │
│  │ • span_id       │          │                              │  │
│  └────────┬────────┘          └─────────────┬────────────────┘  │
└───────────┼──────────────────────────────────┼───────────────────┘
            │                                   │
            │ Logs to stdout/stderr             │ OTLP HTTP (4318)
            │                                   │
            ▼                                   ▼
┌───────────────────────┐         ┌────────────────────────────┐
│   Promtail            │         │   Tempo                    │
│   (Log Collector)     │         │   (Trace Storage)          │
│                       │         │                            │
│ • Reads Docker logs   │         │ • Receives OTLP traces     │
│ • Ships to Loki       │         │ • Stores trace data        │
│ • Labels by container │         │ • Query by trace ID        │
└───────────┬───────────┘         └─────────────┬──────────────┘
            │                                   │
            │ Push logs                         │ Query traces
            │                                   │
            ▼                                   │
┌───────────────────────┐                      │
│   Loki                │◄─────────────────────┘
│   (Log Aggregation)   │  Query logs by trace_id
│                       │
│ • Stores logs         │
│ • Indexes by labels   │
│ • Query by trace_id   │
└───────────┬───────────┘
            │
            │ Query logs & traces
            │
            ▼
┌─────────────────────────────────────────────┐
│   Grafana (Visualization - Port 3000)       │
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │   Explore Tab                         │  │
│  │   • Loki data source                  │  │
│  │   • Tempo data source                 │  │
│  │   • Bi-directional navigation         │  │
│  └──────────────────────────────────────┘  │
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │   Observability Dashboard             │  │
│  │   • Request rate graphs               │  │
│  │   • Log level breakdown               │  │
│  │   • Error tracking                    │  │
│  │   • Live log stream                   │  │
│  │   • Operations table with trace IDs   │  │
│  └──────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
            │
            │ Stores state
            ▼
┌─────────────────────────────────────────────┐
│   MongoDB (Port 27017)                      │
│   • Application data (tasks)                │
│   • Grafana dashboards (if persisted)      │
└─────────────────────────────────────────────┘
```

## Component Details

### 1. Go TODO API

**Location:** Running in Docker container `go-todo-api-api-1`

**Key Features:**
- RESTful API with Huma v2 framework
- Chi router for HTTP routing
- Structured JSON logging with `log/slog`
- OpenTelemetry instrumentation for tracing
- Context propagation throughout request lifecycle

**Logging:**
```go
// Example structured log with trace context
logger.WithTrace(ctx).Info("Created new task",
    "title", task.Title,
    "id", task.ID,
)
// Output:
// {"time":"2025-11-08T20:47:47.833Z","level":"INFO","msg":"Created new task",
//  "trace_id":"1a7cd4fc0152e75a7bbd65d421e58c21","span_id":"70932dc2af870929",
//  "title":"Task title","id":"690fac73..."}
```

**Tracing:**
- Each HTTP request creates a root span
- Database operations create child spans
- Context is propagated through `context.Context`
- Traces exported via OTLP HTTP to Tempo on port 4318

**Configuration:**
- `MONGO_URI`: MongoDB connection string
- `API_KEY`: API authentication key
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Tempo endpoint (defaults to localhost:4318)
- `PORT`: API server port (default 8080)

### 2. Loki (Log Aggregation)

**Location:** `go-todo-api-loki-1`
**Port:** 3100

**Purpose:**
- Aggregates logs from all containers via Promtail
- Indexes logs by labels (container, job, service_name)
- Stores logs efficiently with compression
- Provides LogQL query language

**Key Features:**
- Label-based indexing (not full-text)
- Efficient storage for high-volume logs
- Integrates with Grafana for visualization
- Supports log filtering and parsing

**Example Query:**
```logql
{container="go-todo-api-api-1"} | json | trace_id="1a7cd4fc..."
```

### 3. Tempo (Distributed Tracing)

**Location:** `go-todo-api-tempo-1`
**Ports:**
- 3200: HTTP API for queries
- 4317: OTLP gRPC endpoint
- 4318: OTLP HTTP endpoint (used by API)

**Purpose:**
- Receives traces from API via OpenTelemetry Protocol (OTLP)
- Stores trace data with spans
- Provides trace query API
- Supports TraceQL query language

**Key Features:**
- Low-cost trace storage
- Fast trace lookups by ID
- Service graph generation
- Span-level detail with timing

**Trace Structure:**
```
Trace ID: 1a7cd4fc0152e75a7bbd65d421e58c21
└── Span: POST /tasks (root span)
    ├── Duration: 2.5ms
    ├── Attributes: http.method=POST, http.url=/tasks
    └── Child Span: MongoDB InsertOne
        ├── Duration: 1.2ms
        └── Attributes: db.operation=insert, db.collection=tasks
```

### 4. Promtail (Log Shipper)

**Location:** `go-todo-api-promtail-1`
**Port:** 9080 (metrics)

**Purpose:**
- Collects logs from Docker containers
- Ships logs to Loki
- Adds labels for filtering and organization

**Configuration:** `promtail-config.yaml`
```yaml
# Connects to Docker socket to read container logs
- job_name: docker
  docker_sd_configs:
    - host: unix:///var/run/docker.sock

  # Adds labels from container metadata
  relabel_configs:
    - source_labels: ['__meta_docker_container_name']
      target_label: 'container'
```

**How It Works:**
1. Monitors Docker socket for new log lines
2. Reads logs from all containers
3. Adds labels (container name, job, etc.)
4. Batches and ships logs to Loki
5. Tracks position to avoid duplicate shipping

### 5. Grafana (Visualization)

**Location:** `go-todo-api-grafana-1`
**Port:** 3000

**Purpose:**
- Visualizes logs from Loki
- Visualizes traces from Tempo
- Provides dashboards for metrics
- Enables trace-to-log correlation

**Data Sources:**

**Loki Configuration:**
```yaml
- name: Loki
  type: loki
  url: http://loki:3100
  jsonData:
    derivedFields:
      - datasourceUid: Tempo
        matcherRegex: "trace_id\"?:\"?([0-9a-fA-F]+)"
        name: TraceID
        url: $${__value.raw}
```
- Extracts `trace_id` from logs using regex
- Creates clickable link to Tempo trace
- Enables "Logs → Trace" navigation

**Tempo Configuration:**
```yaml
- name: Tempo
  type: tempo
  url: http://tempo:3200
  jsonData:
    tracesToLogs:
      datasourceUid: Loki
      filterByTraceID: true
      filterBySpanID: false
```
- Configures link from traces to logs
- Filters logs by trace ID automatically
- Enables "Trace → Logs" navigation

**Dashboard:** `go-todo-api` (UID: `go-todo-api`)

**Panels:**
1. **Request Rate** - Time series showing requests/minute
2. **Log Levels Over Time** - Breakdown of INFO/ERROR/WARN logs
3. **Recent API Logs** - Live log stream with full context
4. **Total Requests (Last 5 min)** - Stat panel with thresholds
5. **Error Count (Last 5 min)** - Error tracking with alerts
6. **Recent Operations with Trace IDs** - Table view of operations

### 6. MongoDB (Database)

**Location:** `go-todo-api-mongodb-1`
**Port:** 27017

**Purpose:**
- Application database for TODO tasks
- Stores task documents
- Provides CRUD operations

**Not directly part of observability** but included in the stack for the application to function.

## Data Flow

### Request Flow (Happy Path)

```
1. Client sends HTTP request to API
   ↓
2. API creates root span (OpenTelemetry)
   ↓
3. Logging middleware logs request details with trace_id
   ↓
4. Request handler executes business logic
   ↓
5. Handler creates child span for DB operation
   ↓
6. Handler logs operation result with trace_id
   ↓
7. API exports trace to Tempo via OTLP
   ↓
8. API logs to stdout/stderr
   ↓
9. Promtail reads logs from Docker
   ↓
10. Promtail ships logs to Loki
    ↓
11. User views trace in Grafana Tempo
    ↓
12. User clicks "Related Logs" → sees correlated logs in Loki
```

### Trace-to-Log Correlation Flow

**Scenario 1: Starting from Logs**
```
1. User queries Loki: {container="go-todo-api-api-1"} | json
2. Grafana displays logs with parsed JSON fields
3. User sees trace_id field in log entry
4. User clicks trace_id (Tempo icon appears)
5. Grafana navigates to Tempo with trace_id
6. Full trace displayed with all spans
```

**Scenario 2: Starting from Trace**
```
1. User queries Tempo with trace ID
2. Grafana displays trace waterfall diagram
3. User clicks on a span
4. User clicks "Logs for this span" button
5. Grafana queries Loki with: {container="..."} | json | trace_id="..."
6. Related logs displayed
```

## Key Configuration Files

### docker-compose.yml

Orchestrates all services:
- Defines 6 services: API, MongoDB, Loki, Tempo, Grafana, Promtail
- Sets up Docker network (`lgtm`)
- Mounts configuration files
- Sets environment variables

**Key sections:**
```yaml
api:
  environment:
    - OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318  # Enables trace export
  depends_on:
    - mongodb
    - tempo
    - loki
```

### grafana-datasources.yaml

Configures Grafana data sources with correlation:
- Loki data source with `derivedFields` for Logs → Tempo
- Tempo data source with `tracesToLogs` for Tempo → Loki

### grafana-dashboard.json

Defines the observability dashboard:
- 6 panels with LogQL queries
- Time series, stat, logs, and table visualizations
- Auto-refresh every 5 seconds
- Color-coded thresholds

### promtail-config.yaml

Configures log collection:
- Docker socket connection
- Label extraction from container metadata
- Log shipping to Loki

### tempo-config.yaml

Configures Tempo:
- OTLP receiver endpoints (gRPC and HTTP)
- Storage configuration
- Trace retention policies

## Usage Examples

### Viewing Logs with Trace Context

**In Grafana Explore (Loki):**
```logql
# All API logs with JSON parsing
{container="go-todo-api-api-1"} | json

# Filter by specific trace ID
{container="go-todo-api-api-1"} | json | trace_id="1a7cd4fc..."

# Filter by log level
{container="go-todo-api-api-1"} | json | level="ERROR"

# Find logs about specific task
{container="go-todo-api-api-1"} | json | title=~".*test.*"
```

### Querying Traces

**In Grafana Explore (Tempo):**
```
# Direct trace ID lookup
1a7cd4fc0152e75a7bbd65d421e58c21

# TraceQL query (all traces)
{}

# TraceQL with filters
{ span.http.status_code = 200 }
```

### Viewing the Dashboard

Navigate to: `http://localhost:3000/d/go-todo-api`

**Features:**
- Auto-refreshes every 5 seconds
- Adjustable time range (default: last 15 minutes)
- Click any trace_id in logs to jump to trace
- Color-coded alerts (green/yellow/red)

## Debugging Guide

### Problem: No logs appearing in Loki

**Check:**
1. Is Promtail running? `docker ps | grep promtail`
2. Check Promtail logs: `docker logs go-todo-api-promtail-1`
3. Verify Loki is receiving data: `curl http://localhost:3100/loki/api/v1/labels`
4. Check if logs are being generated: `docker logs go-todo-api-api-1`

**Common causes:**
- Promtail not connected to Docker socket
- Logs generated before Promtail started (only ships new logs)
- Loki ingester not ready (wait 15-20 seconds after startup)

### Problem: No traces in Tempo

**Check:**
1. Is Tempo running? `docker ps | grep tempo`
2. Check API is exporting: `docker logs go-todo-api-api-1 | grep "trace"`
3. Verify OTLP endpoint: Should see no "connection refused" errors
4. Test Tempo API: `curl http://localhost:3200/api/traces/<trace-id>`

**Common causes:**
- `OTEL_EXPORTER_OTLP_ENDPOINT` not set correctly
- Tempo not accessible from API container (network issue)
- Traces exported before Tempo started

### Problem: Correlation not working

**Check:**
1. Verify trace_id in logs: Query Loki and check JSON has `trace_id` field
2. Verify trace exists: Query Tempo with the trace_id from logs
3. Check Grafana datasource config: Should have `derivedFields` (Loki) and `tracesToLogs` (Tempo)
4. Restart Grafana: `docker restart go-todo-api-grafana-1`

**Common causes:**
- Old Grafana configuration cached
- Trace ID format mismatch
- Datasource UIDs incorrect

## Performance Considerations

### Log Volume
- **Current setup:** Logs every request with full JSON context
- **For production:** Consider sampling (log 1% of requests) or filtering (only ERROR level)
- **Storage:** Loki stores logs efficiently, but high volume needs retention policies

### Trace Sampling
- **Current setup:** Samples 100% of traces (`AlwaysSample`)
- **For production:** Use `TraceIdRatioBased(0.1)` to sample 10% of traces
- **Trade-off:** Lower overhead vs. less visibility

### Dashboard Query Load
- **Auto-refresh:** Set to 5 seconds for development
- **For production:** Increase to 30-60 seconds to reduce load
- **Time ranges:** Shorter ranges = faster queries

## Security Considerations

### Secrets Management
- **Never commit:** `.env` file is in `.gitignore`
- **Docker secrets:** Consider using Docker secrets for MongoDB credentials
- **API keys:** Rotate API keys regularly

### Network Isolation
- **Current setup:** All services on same Docker network (`lgtm`)
- **For production:** Use separate networks (frontend, backend, observability)
- **Grafana:** Currently allows anonymous access (development only)

### Log Sanitization
- **PII data:** Be careful logging user data (emails, passwords, etc.)
- **Credentials:** Never log database connection strings or API keys
- **Current code:** Logs task titles and IDs (generally safe)

## Extending the System

### Adding Prometheus Metrics

1. Add Prometheus to `docker-compose.yml`
2. Instrument Go code with `prometheus/client_golang`
3. Add metrics scrape config
4. Add Prometheus data source to Grafana
5. Create metrics dashboards (request latency, throughput, etc.)

### Adding Alerts

1. Configure Grafana alerting in dashboard panels
2. Set up notification channels (email, Slack, PagerDuty)
3. Define alert rules (e.g., error rate > 5%)
4. Test alert firing

### Multiple Environments

1. Use Docker Compose profiles or separate files
2. Add environment selector to Grafana
3. Use different Loki/Tempo instances per environment
4. Add `environment` label to logs and traces

## Resources

### Documentation
- [Grafana Loki](https://grafana.com/docs/loki/)
- [Grafana Tempo](https://grafana.com/docs/tempo/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
- [Promtail](https://grafana.com/docs/loki/latest/send-data/promtail/)

### Dashboards
- **Main Dashboard:** http://localhost:3000/d/go-todo-api
- **Explore Loki:** http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22Loki%22%7D
- **Explore Tempo:** http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22Tempo%22%7D

### API Endpoints
- **API Docs:** http://localhost:8080/docs
- **Loki API:** http://localhost:3100
- **Tempo API:** http://localhost:3200
- **Promtail Metrics:** http://localhost:9080/metrics

## Summary

This observability architecture provides:

✅ **Full visibility** into application behavior
✅ **Correlation** between logs and traces
✅ **Real-time monitoring** via dashboards
✅ **Debugging capability** for production issues
✅ **Scalable design** ready for production use

The LGTM stack (Loki, Grafana, Tempo, MongoDB) offers a complete, open-source observability solution that rivals commercial platforms like Datadog or New Relic, at zero cost.
