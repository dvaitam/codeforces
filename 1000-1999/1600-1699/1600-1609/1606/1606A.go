package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var s string
		fmt.Fscan(reader, &s)
		if len(s) > 0 {
			bs := []byte(s)
			bs[len(bs)-1] = bs[0]
			s = string(bs)
		}
		fmt.Fprintln(writer, s)
	}
}
