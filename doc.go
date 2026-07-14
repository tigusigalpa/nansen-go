// Package nansen provides a fast, concurrent-safe, and idiomatic Go client
// for the Nansen AI API (https://docs.nansen.ai/).
//
// The client has zero third-party dependencies, relying exclusively on the
// Go standard library (net/http, encoding/json, context, time). Every
// service method accepts a context.Context as its first argument to support
// timeouts and cancellation, and the Client (along with its services) is
// safe for concurrent use across multiple goroutines.
//
// Construct a client with New, configuring it with functional options such
// as WithBaseURL, WithHTTPClient, WithTimeout, and WithRetry:
//
//	client, err := nansen.New(os.Getenv("NANSEN_API_KEY"),
//	    nansen.WithTimeout(30*time.Second),
//	    nansen.WithRetry(3, 500*time.Millisecond, 5*time.Second),
//	)
//
// Endpoints are grouped into services exposed as fields on Client:
// SmartMoney, TokenGodMode, Profiler, Portfolio, and Historical.
package nansen
