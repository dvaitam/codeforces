package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x, ans int
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			fmt.Fscan(in, &x)
			if x == 1 {
				ans = abs(i-3) + abs(j-3)
			}
		}
	}

	fmt.Fprint(out, ans)
}
