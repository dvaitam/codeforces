package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	if _, err := fmt.Fscan(in, &N); err != nil {
		return
	}
	for i := 0; i < N; i++ {
		var a, b, c, d, e int
		fmt.Fscan(in, &a, &b, &c, &d, &e)
		idx := a + b*3 + c*9
		ch := 'a' + rune(idx%26)
		out.WriteByte(byte(ch))
	}
	out.WriteByte('\n')
}
