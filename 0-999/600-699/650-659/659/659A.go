package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}
	pos := (a - 1 + b) % n
	if pos < 0 {
		pos += n
	}
	fmt.Println(pos + 1)
}
