package main

import (
	"bufio"
	"fmt"
	"os"
)

func countFizzBuzz(n int64) int64 {
	if n < 0 {
		return 0
	}
	fullBlocks := n / 15
	ans := fullBlocks * 3
	remainders := []int64{1, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
	ans += remainders[n%15]
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, countFizzBuzz(n))
	}
}
