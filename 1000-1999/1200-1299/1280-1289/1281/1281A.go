package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		switch {
		case strings.HasSuffix(s, "po"):
			fmt.Fprintln(writer, "FILIPINO")
		case strings.HasSuffix(s, "desu") || strings.HasSuffix(s, "masu"):
			fmt.Fprintln(writer, "JAPANESE")
		case strings.HasSuffix(s, "mnida"):
			fmt.Fprintln(writer, "KOREAN")
		}
	}
}
