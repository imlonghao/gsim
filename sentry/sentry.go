package sentry

import (
	"github.com/getsentry/sentry-go"
	"log"
	"os"
)

var SENTRY *sentry.Hub

func init() {
	scope := sentry.NewScope()
	client, err := sentry.NewClient(sentry.ClientOptions{Dsn: os.Getenv("DSN")})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	SENTRY = sentry.NewHub(client, scope)
}
