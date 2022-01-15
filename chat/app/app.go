package app

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"parallel_world/chat/chatcore/chatroom"
	"parallel_world/chat/chatcore/common"
	"parallel_world/chat/chatcore/one2one"
	"parallel_world/chat/config"
	"parallel_world/chat/helper"
	"parallel_world/chat/myhost"
	"parallel_world/chat/service/mymdns"
	"parallel_world/chat/ui/TextUI"
)
var rws []*bufio.ReadWriter

type App struct {

}

func NewApp() *App {
	app := App{}

	return &app
}

func (app *App) Run() {
	log.Println("This is a test log entry2")
	//me := common.NewUser()

	help := flag.Bool("help", false, "Display Help")
	//cfg := config.ParseFlags()

	if *help {
		fmt.Printf("Welcome to this parallel world, good journey!")
		fmt.Println("Usage: ")
		fmt.Println("/pworld -nick=[your-nick-name] -room=[a room name]")

		os.Exit(0)
	}
	//fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	// parse some flags to set our nickname and the room to join
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "parallel-world", "name of chat room to join")
	flag.Parse()

	ctx := context.Background()

	h := myhost.MakeChatHost()

	tracer, err := pubsub.NewJSONTracer("/tmp/trace/trace.json")
	if err != nil {
		panic(err)
	}
	ps, err := pubsub.NewGossipSub(ctx, h, pubsub.WithEventTracer(tracer))

	if err != nil {
		panic(err)
	}

	if err := mymdns.SetupDiscovery(h); err != nil {
		log.Fatal("dns service setup error:", err)
		// todo: try others discover-service
	}

	nick := *nickFlag
	if len(nick) == 0 {
		nick = helper.DefaultNick(h.ID())
	}

	// 存储name:pid
	err = common.SavePid(h.ID(), nick)
	if err != nil {
		log.Fatal("save nick name error:", err)
	}

	room := *roomFlag
	common.SaveFavorRoom(nick, room)

	cr, err := chatroom.JoinChatRoom(ctx, ps, h, nick, room)
	if err != nil {
		log.Fatal("join room error:", err)
	}

	ui := TextUI.NewChatUI(cr)

	if err = ui.Run(); err != nil {
		helper.PrintErr("running text UI error: %s", err)
	}
}

func main() {
	help := flag.Bool("help", false, "Display Help")
	cfg := config.ParseFlags()

	if *help {
		fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
		fmt.Printf("Usage: \n   Run './chat-with-mymdns'\nor Run './chat-with-mymdns -myhost [myhost] -port [port] -rendezvous [string] -pid [proto ID]'\n")

		os.Exit(0)
	}
	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	ctx := context.Background()
	r := rand.Reader

	PrvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(PrvKey),
	)

	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), one2one.HandleOne2OneChatStream)
	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())

	peerChan := mymdns.InitMDNS(host, cfg.RendezvousString)

	/*
	pi := <-peerChan // will block until we discover a peer
	fmt.Println("number of peers:", len(peerChan))
	fmt.Println("Found peer:", pi, ", connecting")
	 */

	var peers []peer.AddrInfo

	for pi := range peerChan {
		peers = append(peers, pi)

		if pi.ID.Pretty() == host.ID().Pretty() {
			fmt.Println("ignore self:", pi.ID)
			continue
		}
		fmt.Println("Found peer:", pi, ", connecting")

		if err := host.Connect(ctx, pi); err != nil {
			fmt.Println("Oops, Connection failed:", err)
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := host.NewStream(ctx, pi.ID, protocol.ID(cfg.ProtocolID))
		if err != nil {
			fmt.Println("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			//go chat.WriteData(rw)
			go one2one.ReadData(rw)
			rws = append(rws, rw)
			fmt.Println("Connected to:", pi)
		}
	}
	fmt.Printf("len of rws: %d\n", len(rws))
	go one2one.ReadStdin()
	fmt.Printf("number of peers: %d\n", len(peers))

	select {}
}


