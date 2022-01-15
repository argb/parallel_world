package chatroom

import (
	"context"
	"encoding/json"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-pubsub"
	"log"
	"parallel_world/chat/chatcore/common"
)

const ChatRoomBufSize = 128

type ChatRoom struct {
	Messages chan *ChatMessage
	PrivateMessages chan *ChatMessage // for private chat

	Ctx context.Context
	Ps *pubsub.PubSub
	Stream core.Stream // for private chat
	Topic *pubsub.Topic
	Sub *pubsub.Subscription

	RoomName string
	SelfID peer.ID
	Nick string
	Host host.Host
}

type ChatMessage struct {
	Message string
	SenderID string
	SenderNick string
}

func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, self host.Host, nickname string, roomName string) (*ChatRoom, error) {
	// 创建主题
	topic, err := ps.Join(topicName(roomName))
	if err != nil {
		return nil, err
	}
	// 订阅该主题, 返回一个订阅
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	cr := &ChatRoom{
		Ctx: ctx,
		Ps: ps,
		Topic: topic,
		Sub: sub,
		Host: self,
		SelfID: self.ID(),
		Nick: nickname,
		RoomName: roomName,
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
		PrivateMessages: make(chan *ChatMessage, ChatRoomBufSize),
	}

	go cr.readLoop()
	//go cr.readLoopP()

	return cr, nil
}

func topicName(roomName string) string {
	return "chat-room:" + roomName
}

func (cr *ChatRoom) Publish(message string) error {
	m := ChatMessage{
		Message: message,
		SenderID: cr.SelfID.Pretty(),
		SenderNick: cr.Nick,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.Topic.Publish(cr.Ctx, msgBytes)
}

func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.Sub.Next(cr.Ctx)
		if err != nil {
			close(cr.Messages)
			return
		}

		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.SelfID {
			continue
		}
		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			log.Printf("invalid msg: %v\n", msg)
			continue
		}
		cr.Messages <- cm
	}
}

func (cr *ChatRoom) ReadLoopP() {
	for {
		var buffer []byte
		_, err := cr.Stream.Read(buffer)
		if err != nil {
			close(cr.Messages)
			return
		}

		cm := new(ChatMessage)
		err = json.Unmarshal(buffer, cm)
		if err != nil {
			log.Printf("invalid msg: %v\n", buffer)
			continue
		}
		err = common.SavePid(peer.ID(cm.SenderID), cm.SenderNick)
		if err != nil {
			log.Fatal("save pid error:", err)
		}
		cr.PrivateMessages <- cm
	}
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.Ps.ListPeers(topicName(cr.RoomName))
}
