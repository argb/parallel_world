package cmd

import (
	"log"
	"strings"
)

const (
	AtPrefix = "@" // cmd example : @wg, @trump
	Separator = ":" // input example: @wg: hello, good morning
	SlashPrefix = "/"
)
const (
	Quit = "/quit"
)

type CMD interface {

}

type AtCmd struct {
	Operator string
	Operand string
}

type QuitCmd struct {
	Operator string
}
/*
// 处理交互过程中输入的指令

func ReadFromStdin() {
	var reader *bufio.Reader
	reader = bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("read stdin error:", err)
		}

		cmd, data := ParseInput(input)

		if !isValidCmd(cmd) {
			one2more.SendTextToAll(input)
		}else {
			uname := strings.TrimLeft(cmd, cmdAt)
			one2one.SendTextTo(uname, data)
		}

	}

}

 */

func ParseInput(input string) (CMD, string) {
	// todo 好好测试下
	inputs := strings.SplitN(input, Separator, 2)
	if len(inputs) == 2 {
		cmd := inputs[0]
		data := inputs[1]
		if strings.HasPrefix(cmd, AtPrefix) {
			if len(cmd) < 2 {
				log.Printf("<%s> is not a valid cmd", cmd)
				return nil, input // invalid cmd, means no cmd exist, so the data is whole input
			}else {
				atCmd := AtCmd{
					Operator: AtPrefix,
					Operand: strings.TrimLeft(cmd, AtPrefix),
				}
				return atCmd, data
			}
		}
	}else {
		if strings.HasPrefix(input, SlashPrefix) {
			switch input {
			case Quit:
				cmd := QuitCmd{
				Operator: Quit,
				}
				return cmd, ""
			}
		}
	}

	return nil, input
}
