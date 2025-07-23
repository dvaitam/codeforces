package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	pos1, posn := -1, -1
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x == 1 {
			pos1 = i
		}
		if x == n {
			posn = i
		}
	}

	if pos1 == -1 || posn == -1 {
		return
	}

	dist := func(a, b int) int {
		if a > b {
			return a - b
		}
		return b - a
	}

	ans := dist(pos1, 0)
	if v := dist(pos1, n-1); v > ans {
		ans = v
	}
	if v := dist(posn, 0); v > ans {
		ans = v
	}
	if v := dist(posn, n-1); v > ans {
		ans = v
	}

	fmt.Fprintln(writer, ans)
}
