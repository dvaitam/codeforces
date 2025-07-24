package main

import (
	"bufio"
	"fmt"
	"os"
)

var memo [105][105][2][2]int8

func dfs(e, o, turn, parity int) bool {
	if e == 0 && o == 0 {
		return parity == 0
	}
	m := &memo[e][o][turn][parity]
	if *m != 0 {
		return *m == 1
	}
	var res bool
	if turn == 0 {
		// Alice's turn: she tries to find any move leading to a win
		if e > 0 && dfs(e-1, o, 1, parity) {
			res = true
		} else if o > 0 && dfs(e, o-1, 1, parity^1) {
			res = true
		} else {
			res = false
		}
	} else {
		// Bob's turn: he tries to force Alice's loss
		res = true
		if e > 0 && !dfs(e-1, o, 0, parity) {
			res = false
		}
		if res && o > 0 && !dfs(e, o-1, 0, parity) {
			res = false
		}
	}
	if res {
		*m = 1
	} else {
		*m = 2
	}
	return res
}

func solve(e, o int) string {
	// reset memo for each test case
	for i := 0; i <= e; i++ {
		for j := 0; j <= o; j++ {
			memo[i][j][0][0] = 0
			memo[i][j][0][1] = 0
			memo[i][j][1][0] = 0
			memo[i][j][1][1] = 0
		}
	}
	if dfs(e, o, 0, 0) {
		return "Alice"
	}
	return "Bob"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		e, o := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x%2 == 0 {
				e++
			} else {
				o++
			}
		}
		fmt.Fprintln(writer, solve(e, o))
	}
}
