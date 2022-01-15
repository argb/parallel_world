package config

import "flag"

const ChatProtocol = "/parallel_world/chat/1.0.0"

type Config struct {
	RendezvousString string
	ProtocolID string
	ListenHost string
	ListenPort int
}

func ParseFlags() *Config {
	c := new(Config)

	flag.StringVar(&c.RendezvousString, "rendezvous", "moyu", "\"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.ProtocolID, "pid", ChatProtocol, "Sets a protocol id for stream headers")
	flag.StringVar(&c.ListenHost, "myhost", "0.0.0.0", "The bootstrap node myhost listen address\n")
	flag.IntVar(&c.ListenPort, "port", 0, "node listen port")

	flag.Parse()

	return c
}
