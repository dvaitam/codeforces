package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	total := 0
	for i := 0; i < n; i++ {
		var w int
		fmt.Fscan(in, &w)
		total += (w + k - 1) / k
	}
	fmt.Println((total + 1) / 2)
}
