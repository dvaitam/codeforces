package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// sum of 1..n
	total := int64(n) * int64(n+1) / 2
	for i := 0; i < n-1; i++ {
		var x int
		fmt.Fscan(reader, &x)
		total -= int64(x)
	}
	fmt.Println(total)
}
