package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	fmt.Fscan(in, &n)
	fmt.Println(n/2 + 1)
}
