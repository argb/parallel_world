package mymdns

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
)

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "parallel_world_chat"

type DiscoveryNotifee struct {
	PeerChan chan peer.AddrInfo
	Host host.Host

	AutoConnect bool
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("discovered new peer %s\n", pi.ID.Pretty())

	//n.PeerChan <- pi

	if n.AutoConnect {
		err := n.Host.Connect(context.Background(), pi)
		if err != nil {
			log.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
		}
		log.Printf("conntected to %v\n", pi.ID.Pretty())
	}
	for _, pa := range pi.Addrs {
		log.Printf("address: %v\n", pa.String())
	}
	log.Printf("addresses %#v saved\n", pi.Addrs)
	n.Host.Peerstore().AddAddrs(pi.ID, pi.Addrs, peerstore.PermanentAddrTTL)
}

func InitMDNS(peerHost host.Host, rendezvous string) chan peer.AddrInfo{
	notifee := new(DiscoveryNotifee)
	notifee.PeerChan = make(chan peer.AddrInfo)
	notifee.Host = peerHost
	notifee.AutoConnect = true

	ser := mdns.NewMdnsService(peerHost, rendezvous, notifee)
	if err := ser.Start(); err != nil {
		panic(err)
	}

	return notifee.PeerChan
}


// SetupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func SetupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	InitMDNS(h, DiscoveryServiceTag)


	return nil
}