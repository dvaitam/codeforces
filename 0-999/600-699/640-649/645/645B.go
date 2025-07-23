package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k >= n/2 {
		// Enough swaps to reverse the entire arrangement
		fmt.Println(n * (n - 1) / 2)
		return
	}
	// Optimal strategy is to swap the first k and last k cows
	ans := k * (2*n - 2*k - 1)
	fmt.Println(ans)
}
