package helper

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
)

func DefaultNick(p peer.ID) string {
	return fmt.Sprintf("%s-%s", os.Getenv("USER"), ShortID(p))
}


func ShortID(pid peer.ID) string {
	pretty := pid.Pretty()
	return pretty[len(pretty) - 8 : ]
}


// PrintErr is like fmt.Printf, but writes to stderr.
func PrintErr(m string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, m, args...)
}