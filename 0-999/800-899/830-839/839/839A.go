package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	candies := 0
	total := 0
	for i := 1; i <= n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		candies += a
		give := 8
		if candies < 8 {
			give = candies
		}
		total += give
		candies -= give
		if total >= k {
			fmt.Println(i)
			return
		}
	}
	fmt.Println(-1)
}
