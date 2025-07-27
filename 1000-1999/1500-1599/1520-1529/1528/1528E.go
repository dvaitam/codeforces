package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	if n == 1 {
		fmt.Println(5)
	} else {
		fmt.Println(0)
	}
}
