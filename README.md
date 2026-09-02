# nansen-go

![Nansen AI Golang SDK](https://i.postimg.cc/hGsGsnDF/nansen-ai-api-golang.jpg)

[![CI](https://github.com/tigusigalpa/nansen-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/nansen-go/actions/workflows/ci.yml)
[![Tests](https://github.com/tigusigalpa/nansen-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/nansen-go/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![CodeQL](https://github.com/tigusigalpa/nansen-go/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/nansen-go/actions/workflows/codeql.yml)
[![Codecov](https://codecov.io/gh/tigusigalpa/nansen-go/graph/badge.svg)](https://codecov.io/gh/tigusigalpa/nansen-go)
[![GitHub Release](https://img.shields.io/github/v/release/tigusigalpa/nansen-go?style=flat-square)](https://github.com/tigusigalpa/nansen-go/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/tigusigalpa/nansen-go)

A Go client for the [Nansen AI API](https://docs.nansen.ai/) that tries to stay out of your way. No third-party
dependencies, no surprises — just the standard library and an API that feels like the rest of your Go code.

**Package:** [pkg.go.dev/github.com/tigusigalpa/nansen-go](https://pkg.go.dev/github.com/tigusigalpa/nansen-go)

> 📖 **[Full documentation available on Wiki](https://github.com/tigusigalpa/nansen-go/wiki)**

## Why you might like it

- **Nothing to vendor.** The whole thing is built on the standard library. `go get` it and you're done — no dependency
  tree to audit.
- **Contexts everywhere.** Every call takes a `context.Context` first, so timeouts and cancellation work exactly the way
  you'd expect.
- **Configured with options, not structs.** `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithRetry` — mix and match
  what you need.
- **Safe to share.** Create one client and hand it to as many goroutines as you like. No locks, no fuss.
- **Typed on purpose.** Chains, sort fields, trader types, and labels are real constants. Optional request fields are
  pointers, so a stray `false` or `0` never sneaks into your JSON.
- **Errors you can actually inspect.** Failures come back as an `*APIError` with the status code, message, raw body, and
  rate-limit headers — and they play nicely with `errors.Is`.
- **Retries when you want them.** Opt in with `WithRetry` and the client backs off exponentially on 429s, honoring
  `Retry-After`, `RateLimit-Reset`, and `X-RateLimit-Reset` along the way.

## Installation

```bash
go get github.com/tigusigalpa/nansen-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/tigusigalpa/nansen-go"
)

func main() {
    client, err := nansen.New(os.Getenv("NANSEN_API_KEY"),
        nansen.WithTimeout(30*time.Second),
        nansen.WithRetry(3, 500*time.Millisecond, 5*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    tf := nansen.Timeframe24H
    resp, err := client.TokenGodMode.TokenScreener(ctx, &nansen.TokenScreenerRequest{
        Chains: []nansen.Chain{nansen.ChainEthereum, nansen.ChainSolana},
        Timeframe: &tf,
        Pagination: &nansen.PaginationRequest{
            Page:    nansen.IntPtr(1),
            PerPage: nansen.IntPtr(20),
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, t := range resp.Data {
        fmt.Printf("%s on %s\n", t.TokenSymbol, t.Chain)
    }
}
```

You'll need a `NANSEN_API_KEY` — grab one from [Nansen](https://nansen.ai/).

## Configuring the client

Everything is optional. Pass only the options you care about; the rest fall back to sensible defaults.

```go
client, err := nansen.New(apiKey,
    nansen.WithBaseURL("https://api.nansen.ai"),
    nansen.WithHTTPClient(&http.Client{Timeout: 10*time.Second}),
    nansen.WithTimeout(30*time.Second),
    nansen.WithRetry(3, 500*time.Millisecond, 5*time.Second),
)
```

| Option           | Description                                                                                            |
|------------------|--------------------------------------------------------------------------------------------------------|
| `WithBaseURL`    | Override the default base URL (`https://api.nansen.ai`).                                               |
| `WithHTTPClient` | Provide a custom `*http.Client`.                                                                       |
| `WithTimeout`    | A default request timeout, used when your context doesn't already carry a deadline.                    |
| `WithRetry`      | Turns on automatic retries with exponential backoff. Takes max attempts, initial delay, and max delay. |

## What's available

The API surface is split into services that hang off the client:

- `client.SmartMoney` — netflows, holdings, and DEX trades.
- `client.TokenGodMode` — the token screener, flow intelligence, and who-bought-sold.
- `client.Profiler` — current balances and DEX trade history for an address.
- `client.Portfolio` — DeFi holdings.
- `client.Historical` — the backtesting endpoints under `/api/v1beta1/`.

## When things go wrong

Every API failure comes back as an `*nansen.APIError`, so you can dig into exactly what happened:

```go
resp, err := client.SmartMoney.Netflow(ctx, req)
if err != nil {
    var apiErr *nansen.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.StatusCode)
        fmt.Println(apiErr.Message)
        fmt.Println(apiErr.RawBody)
        fmt.Println(apiErr.Headers.Get("RateLimit-Remaining"))
        if apiErr.RetryAfter != nil {
            fmt.Println("retry after:", *apiErr.RetryAfter)
        }
        if apiErr.RateLimitRemaining != nil {
            fmt.Println("requests remaining:", *apiErr.RateLimitRemaining)
        }
    }
}
```

If you just want to branch on the kind of failure, reach for `errors.Is` and the built-in sentinels:

```go
if errors.Is(err, nansen.ErrRateLimited) { /* ... */ }
if errors.Is(err, nansen.ErrUnauthorized) { /* ... */ }
if errors.Is(err, nansen.ErrNotFound)      { /* ... */ }
```

## About retries

Once you've called `WithRetry`, the client quietly retries a few situations for you:

- **429 Too Many Requests** — it waits according to `Retry-After` or `RateLimit-Reset`/`X-RateLimit-Reset`, and stashes
  `RateLimit-Remaining`/`X-RateLimit-Remaining` on the returned `APIError` in case you want to peek at your remaining
  quota.
- **Transient 5xx responses** — the kind that usually clear up on a second try.
- **Flaky network errors** — unless your context has already been cancelled.

Backoff grows exponentially but never exceeds the max delay you set. One thing worth knowing: your timeout budget (from
`WithTimeout` or a deadline on the context) covers **all** the attempts together, not each one on its own — so you
always stay within the bound you asked for.

## Optional fields and pointer helpers

Optional request fields are pointers, which means `false`, `0`, and empty strings only get sent when you actually mean
them:

```go
req := &nansen.TokenScreenerRequest{
    Chains: []nansen.Chain{nansen.ChainSolana},
    Filters: &nansen.TokenScreenerFilters{
        IncludeStablecoins: nansen.BoolPtr(false),
        MarketCapUSD: &nansen.NumericRangeFilter{
            Min: nansen.Float64Ptr(1_000_000),
            Max: nansen.Float64Ptr(50_000_000),
        },
    },
}
```

To keep that from getting tedious, there are little helpers for the common types:

- `StringPtr(s string) *string`
- `IntPtr(i int) *int`
- `BoolPtr(b bool) *bool`
- `Float64Ptr(f float64) *float64`

## Examples

If you'd rather learn by running something, the [`examples`](examples/) directory has a few complete programs:

- [`examples/screener`](examples/screener/main.go) — the token screener with filters and sorting.
- [`examples/profiler`](examples/profiler/main.go) — address balances and DEX trade history.
- [`examples/smart_money`](examples/smart_money/main.go) — Smart Money netflows, holdings, and DEX trades.

Point one at your API key and go:

```bash
NANSEN_API_KEY=your_api_key go run ./examples/screener
```

## Endpoints covered

### Smart Money

- `POST /api/v1/smart-money/netflow`
- `POST /api/v1/smart-money/holdings`
- `POST /api/v1/smart-money/dex-trades`

### Token God Mode & Screener

- `POST /api/v1/token-screener`
- `POST /api/v1/tgm/flow-intelligence`
- `POST /api/v1/tgm/who-bought-sold`

### Profiler

- `POST /api/v1/profiler/address/current-balance`
- `POST /api/v1/profiler/dex-trades`

### Portfolio

- `POST /api/v1/portfolio/defi-holdings`

### Historical Data (Backtesting)

- `POST /api/v1beta1/tgm/historical-token-flow-summary`
- `POST /api/v1beta1/smart-money/historical-token-balances`

## Testing

There's a unit test suite covering the retry/backoff logic, error mapping, and option validation. Run it the usual way:

```bash
go test ./...
```

## License

MIT © Igor Sazonov. See [LICENSE](LICENSE).
