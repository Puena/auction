package nats

import "github.com/nats-io/nats.go"

type Config struct {
	AppName string
	URL     string
}

func Connect(conf Config) (*nats.Conn, error) {
	return nats.Connect(conf.URL, nats.Name(conf.AppName))
}
