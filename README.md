# Goverland Core WEB API

<a href="https://github.com/goverland-labs/goverland-core-web-api?tab=License-1-ov-file" rel="nofollow"><img src="https://img.shields.io/github/license/goverland-labs/goverland-core-web-api" alt="GPL 3.0" style="max-width:100%;"></a>
![unit-tests](https://github.com/goverland-labs/goverland-core-web-api/workflows/unit-tests/badge.svg)
![golangci-lint](https://github.com/goverland-labs/goverland-core-web-api/workflows/golangci-lint/badge.svg)

Public REST API gateway for the [Goverland](https://goverland.xyz) platform. It exposes DAO governance data over HTTP by proxying requests to internal gRPC services ([core-storage](https://github.com/goverland-labs/goverland-core-storage), core-feed).

## API Endpoints

All endpoints are prefixed with `/v1` (configurable).

### DAOs
| Method | Path | Description |
|---|---|---|
| GET/POST | `/daos` | List DAOs (with filtering) |
| GET | `/daos/top` | Top DAOs by category |
| GET | `/daos/recommendations` | DAO recommendations |
| GET | `/daos/{id}` | Get DAO by ID |
| GET | `/daos/{id}/feed` | DAO feed |
| GET | `/daos/{id}/delegates` | DAO delegates |
| GET | `/daos/{id}/delegate-profile` | Delegate profile |
| GET | `/daos/{id}/token-info` | Token info |
| GET | `/daos/{id}/token-chart` | Token price chart |

### Proposals
| Method | Path | Description |
|---|---|---|
| GET | `/proposals` | List proposals (with filtering and sorting) |
| GET | `/proposals/top` | Top proposals |
| GET | `/proposals/{id}` | Get proposal by ID |
| GET | `/proposals/{id}/votes` | Proposal votes |
| POST | `/proposals/{id}/votes/validate` | Validate vote |
| POST | `/proposals/{id}/votes/prepare` | Prepare vote |
| POST | `/proposals/votes` | Submit vote |

### Users & Delegates
| Method | Path | Description |
|---|---|---|
| GET | `/user/{address}/votes` | User vote history |
| GET | `/user/{address}/participated-daos` | DAOs user voted in |
| GET | `/user/{address}/delegates/top` | User's top delegates |
| GET | `/user/{address}/delegators/top` | User's top delegators |
| GET | `/user/{address}/delegations/total` | Delegation summary |

### Other
| Method | Path | Description |
|---|---|---|
| GET | `/ens-name` | Resolve addresses to ENS names |
| GET | `/ens-address` | Resolve ENS names to addresses |
| GET | `/stats/totals` | Platform statistics |
| POST | `/feed` | Get feed by filters |
| POST | `/subscribe` | Create subscriber |
| PUT | `/subscribe` | Update subscriber |
| POST | `/subscriptions` | Subscribe to DAO |
| DELETE | `/subscriptions` | Unsubscribe from DAO |

## Project Structure

```
main.go                      # Entry point
internal/
  app.go                     # Application bootstrap
  config/                    # Environment-based configuration
  grpc/                      # gRPC client connections to backend services
  rest/
    handlers/                # REST endpoint handlers
    form/                    # Request parsing and validation
    models/                  # Response models (dao, proposal, delegate, etc.)
    helpers/                 # Shared REST utilities
  response/                  # Response formatting and error handling
  helpers/                   # General helpers
  logger/                    # Zerolog logger setup
pkg/
  grpcsrv/                   # gRPC server utilities
  health/                    # Health check endpoint
  helpers/                   # Shared helpers
  middleware/                # HTTP/gRPC middleware
  prometheus/                # Prometheus metrics handler
```

## Build & Run

```bash
go build ./...
go test ./...
golangci-lint run
```

## Configuration

Environment variables (parsed with [caarlos0/env](https://github.com/caarlos0/env)):

| Variable | Description |
|---|---|
| `LOG_LEVEL` | Zerolog level (default: `info`) |
| `REST_LISTEN` | REST API listen address (default: `:8080`) |
| `REST_API_VERSION` | API version prefix (default: `v1`) |
| `INTERNAL_API_CORE_STORAGE_ADDRESS` | Core storage gRPC address (default: `localhost:11100`) |
| `INTERNAL_API_CORE_FEED_ADDRESS` | Core feed gRPC address (default: `localhost:11000`) |
| `PROMETHEUS_LISTEN` | Prometheus metrics address (default: `:2112`) |
| `HEALTH_LISTEN` | Health check address (default: `:3000`) |

## Contribution Rules

[CONTRIBUTING.md](CONTRIBUTING.md)

## Changelog

[CHANGELOG.md](CHANGELOG.md)
