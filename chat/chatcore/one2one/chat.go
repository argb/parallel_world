package one2one

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
	"os"
	"parallel_world/chat/chatcore/common"
)

type Chat struct {

}
var UserInfo map[string]peer.ID

var rws []*bufio.ReadWriter

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

func HandleOne2OneChatStream(stream network.Stream) {
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