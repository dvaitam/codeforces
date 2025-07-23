package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the actual algorithm. For now it prints 0.

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	var p int64
	if _, err := fmt.Fscan(reader, &n, &k, &p); err != nil {
		return
	}
	fmt.Println(0)
}
