package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	in.ReadString('\n')
	for ; t > 0; t-- {
		line, _ := in.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		b := []byte(line)
		if len(b) >= 5 {
			b[0], b[4] = b[4], b[0]
		}
		fmt.Fprintln(out, string(b))
	}
}
