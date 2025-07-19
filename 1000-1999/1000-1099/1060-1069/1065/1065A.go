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

	var t int64
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := int64(0); i < t; i++ {
		var s, a, b, c int64
		fmt.Fscan(reader, &s, &a, &b, &c)
		x := s / c
		res := x + (x/a)*b
		fmt.Fprintln(writer, res)
	}
}
