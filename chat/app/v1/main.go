package main

import (
	"context"
	"flag"
	"github.com/libp2p/go-libp2p-pubsub"
	"log"
	"parallel_world/chat/chatcore/chatroom"
	"parallel_world/chat/helper"
	"parallel_world/chat/myhost"
	"parallel_world/chat/service/mymdns"
	"parallel_world/chat/ui/TextUI"
	"time"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour


func main() {
	// parse some flags to set our nickname and the room to join
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "awesome-chat-room", "name of chat room to join")
	flag.Parse()

	ctx := context.Background()

	h := myhost.MakeChatHost()

	ps, err := pubsub.NewGossipSub(ctx, h)
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

	room := *roomFlag

	cr, err := chatroom.JoinChatRoom(ctx, ps, h, nick, room)
	if err != nil {
		log.Fatal("join room error:", err)
	}

	ui := TextUI.NewChatUI(cr)

	if err = ui.Run(); err != nil {
		helper.PrintErr("running text UI error: %s", err)
	}

}