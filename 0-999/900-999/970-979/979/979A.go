package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n == 0 {
		fmt.Println(0)
		return
	}
	m := n + 1
	if m%2 == 0 {
		fmt.Println(m / 2)
	} else {
		fmt.Println(m)
	}
}
