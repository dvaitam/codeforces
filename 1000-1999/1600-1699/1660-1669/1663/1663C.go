package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	sum := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		sum += x
	}
	fmt.Println(sum)
}
