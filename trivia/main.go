package main

import (
	"fmt"
	"strings"
)

func main() {
	t1 := "/quit"
	t2 := "@wg: hello"

	inputs := strings.SplitN(t1, ":", 2)
	fmt.Printf("inputs: %#v\n", inputs)

	inputs = strings.SplitN(t2, ":", 2)
	fmt.Printf("inputs: %#v\n", inputs)
	var ids []int
	ids =make([]int, 0)
	ids = append(ids, 10)
	fmt.Println(ids)
}
