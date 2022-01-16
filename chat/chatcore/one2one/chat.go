package one2one

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
	"os"
	"parallel_world/chat/chatcore/common"
)

const ChatBufSize = 128
var PrivateMessages chan *common.ChatMessage

type Chat struct {
	Stream network.Stream
	PrivateMessages chan *common.ChatMessage // for private chat

	MyID peer.ID
	MyNick string

}

var UserInfo map[string]peer.ID
var CurrentChat *Chat

var rws []*bufio.ReadWriter

func init() {
	PrivateMessages = make(chan *common.ChatMessage, ChatBufSize)
}

func NewOne2OneChat(stream network.Stream, myId peer.ID, nick string) *Chat {
	chat := &Chat{
		PrivateMessages: PrivateMessages,
		MyID: myId,
		MyNick: nick,
	}
	CurrentChat = chat

	return chat
}

func HandleOne2OneChatStream1(stream network.Stream) {
	fmt.Println("Got a new stream!")
	remoteID := stream.Conn().RemotePeer().String()
	err := common.SavePid(peer.ID(remoteID), "wg")
	if err != nil {
		log.Fatal("save pid error:", err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	rws = append(rws, rw)
	fmt.Println("number of streams:", len(rws))
	go ReadStdin()

	go ReadData(rw)
	//go WriteData(rw)

}

func HandleOne2OneChatStream(stream network.Stream)  {
	pid, err := peer.Decode(stream.ID())
	if err != nil {
		log.Fatal("pid decode error:", err)
	}

	chat := NewOne2OneChat(stream, pid, common.MyNick)

	go chat.ReadLoop()
}

func (chat *Chat) ReadLoop() {
	stream := chat.Stream
	for {
		var buffer []byte
		_, err := stream.Read(buffer)
		if err != nil {
			close(PrivateMessages)
			return
		}

		cm := new(common.ChatMessage)
		err = json.Unmarshal(buffer, cm)
		if err != nil {
			log.Printf("invalid msg: %v\n", buffer)
			continue
		}
		err = common.SavePid(peer.ID(cm.SenderID), cm.SenderNick)
		if err != nil {
			log.Fatal("save pid error:", err)
		}
		PrivateMessages <- cm
	}
}

func (chat *Chat) SendMassage(message string) {
	m := common.ChatMessage{
		Message: message,
		SenderID: chat.MyID.Pretty(),
		SenderNick: chat.MyNick,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		log.Fatal("encode msg error:", err)
	}
	n, err := chat.Stream.Write(msgBytes)
	log.Printf("%d bytes sended\n", n)
	if err != nil {
		log.Fatal("send chat msg error: ", err)
	}
	if n < len(msgBytes) {
		log.Println("lost msg? -->", message)
	}
}

func ReadStdin() {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from stdin")
			panic(err)
		}
		for i, rw := range rws {
			fmt.Println("send:", i)
			go WriteData(rw, sendData)
			//go ReadData(rw)
		}
	}
}

func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func WriteData(rw *bufio.ReadWriter, data string) {
	//stdReader := bufio.NewReader(os.Stdin)

	//for {
		/*
		fmt.Print("W> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from stdin")
			panic(err)
		}

		 */

		_, err := rw.WriteString(fmt.Sprintf("%s\n", data))
		if err != nil {
			fmt.Println("error writing to buffer")
			panic(err)
		}
		if err = rw.Flush(); err != nil {
			log.Fatalln("error flushing buffer", err)
		}
	//}
}

func bindPeerIDtoUsername(id peer.ID, uname string) {
	UserInfo[uname] = id
}

func getPIDByUsername(uname string) peer.ID {

	return UserInfo[uname]
}

func SendTextTo(uname string, data string) {

}
