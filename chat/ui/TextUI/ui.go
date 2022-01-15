package TextUI

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/rivo/tview"
	"io"
	"log"
	"parallel_world/chat/chatcore/chatroom"
	"parallel_world/chat/chatcore/common"
	"parallel_world/chat/cmd"
	"parallel_world/chat/helper"
	"parallel_world/chat/protocol"
	"time"
)


// 基于命令行文本界面的UI(Text User Interface, TUI)

type ChatUI struct {
	cr *chatroom.ChatRoom
	app *tview.Application
	peersList *tview.TextView

	msgW io.Writer // 聊天室
	msgWP io.Writer // 私聊
	inputChan chan string
	inputChanP chan string // 私聊
	doneChan chan struct{}
}

func NewChatUI(cr *chatroom.ChatRoom) *ChatUI {
	app := tview.NewApplication()

	// 私聊窗口
	msgBoxPrivate := tview.NewTextView()
	msgBoxPrivate.SetDynamicColors(true)
	msgBoxPrivate.SetBorder(true)
	msgBoxPrivate.SetTitle(fmt.Sprintf("Room: %s", cr.RoomName))

	msgBoxPrivate.SetChangedFunc(func() {
		app.Draw()
	})


	// 聊天室窗口
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Room: %s", cr.RoomName))

	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	inputChan := make(chan string, 32)
	inputChanP := make(chan string, 32)

	inputField := tview.NewInputField().
		SetLabel(cr.Nick + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := inputField.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}
		command, data := cmd.ParseInput(line)
		switch command.(type) {
		case cmd.AtCmd:
			uname := command.(cmd.AtCmd).Operator
			fmt.Printf("target:%s\n", uname)
			pid := common.GetPIDByUsername(uname)
			s, err := cr.Host.NewStream(cr.Ctx, peer.ID(pid), protocol.ChatOne2OneProtocol)
			if err != nil {
				log.Fatal(err)
			}
			cr.Stream = s
			inputChanP <- data
			go cr.ReadLoopP()
		case cmd.QuitCmd:
			app.Stop()
			return
		default:
			// send the line onto the input chan and reset the field text
			inputChan <- line
		}

		inputField.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")
	peersList.SetChangedFunc(func() { app.Draw() })

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBoxPrivate, 0, 1, false).
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(inputField, 1, 1, true)

	app.SetRoot(flex, true)

	return &ChatUI{
		cr:        cr,
		app:       app,
		peersList: peersList,
		msgW:      msgBox,
		msgWP:	   msgBoxPrivate,
		inputChan:   inputChan,
		doneChan:    make(chan struct{}, 1),
	}
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *ChatUI) Run() error {
	go ui.handleEvents()
	defer ui.end()

	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *ChatUI) end() {
	ui.doneChan <- struct{}{}
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {
	pids := ui.cr.ListPeers()
	//fmt.Println("peers:", pids)
	// clear is not threadsafe so we need to take the lock.
	ui.peersList.Lock()
	ui.peersList.Clear()
	ui.peersList.Unlock()

	for _, pid := range pids {
		fmt.Fprintln(ui.peersList, helper.ShortID(pid))
	}

	ui.app.Draw()
}

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
func (ui *ChatUI) displayChatMessage(cm *chatroom.ChatMessage) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.SenderNick))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, cm.Message)
}

// displaySelfMessage writes a message from ourself to the message window,
// with our nick highlighted in yellow.
func (ui *ChatUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.cr.Nick))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
func (ui *ChatUI) displayPrivateChatMessage(cm *chatroom.ChatMessage) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.SenderNick))
	fmt.Fprintf(ui.msgWP, "%s %s\n", prompt, cm.Message)
}

// displaySelfMessage writes a message from ourself to the message window,
// with our nick highlighted in yellow.
func (ui *ChatUI) displayPrivateSelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.cr.Nick))
	fmt.Fprintf(ui.msgWP, "%s %s\n", prompt, msg)
}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.
func (ui *ChatUI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()
	fmt.Println("start handle event")

	for {
		select {
		case inputP := <- ui.inputChanP:
			_, err := ui.cr.Stream.Write([]byte(inputP))
			if err != nil {
				helper.PrintErr("send private msg error: %s", err)
			}
			ui.displayPrivateSelfMessage(inputP)
		case pm := <-ui.cr.PrivateMessages:
			// when we receive a message from the chat room, print it to the message window
			ui.displayPrivateChatMessage(pm)
		case input := <-ui.inputChan:
			// when the user types in a line, publish it to the chat room and print to the message window
			err := ui.cr.Publish(input)
			if err != nil {
				helper.PrintErr("publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		case m := <-ui.cr.Messages:
			// when we receive a message from the chat room, print it to the message window
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			// refresh the list of peers in the chat room periodically
			ui.refreshPeers()

		case <-ui.cr.Ctx.Done():
			return

		case <-ui.doneChan:
			return
		}
	}
}

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
