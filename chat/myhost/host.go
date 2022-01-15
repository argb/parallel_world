package myhost

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"log"
)

func MakeBasicHost() *host.Host {
	basicHost, err := libp2p.New()
	if err != nil {
		log.Fatal("making myhost error:", err)
	}


	return &basicHost
}

func MakeChatHost() host.Host {
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	return h
}
