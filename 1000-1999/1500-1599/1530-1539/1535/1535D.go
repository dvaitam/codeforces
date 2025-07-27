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

	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	n := (1 << k) - 1

	// mapping from position in the input string to node index in the tree
	posToNode := make([]int, n+1)
	pos := 1
	for level := k - 1; level >= 0; level-- {
		for node := 1 << level; node <= (1<<(level+1))-1; node++ {
			posToNode[pos] = node
			pos++
		}
	}

	nodeChar := make([]byte, n+1)
	for p := 1; p <= n; p++ {
		node := posToNode[p]
		nodeChar[node] = s[p-1]
	}

	dp := make([]int, n+1)
	var recalc func(int)
	recalc = func(node int) {
		for {
			ch := nodeChar[node]
			left := node * 2
			if left > n {
				if ch == '?' {
					dp[node] = 2
				} else {
					dp[node] = 1
				}
			} else {
				right := left + 1
				if ch == '?' {
					dp[node] = dp[left] + dp[right]
				} else if ch == '0' {
					dp[node] = dp[left]
				} else {
					dp[node] = dp[right]
				}
			}
			if node == 1 {
				break
			}
			node /= 2
		}
	}

	for node := n; node >= 1; node-- {
		recalc(node)
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var p int
		var c string
		fmt.Fscan(reader, &p, &c)
		node := posToNode[p]
		nodeChar[node] = c[0]
		recalc(node)
		fmt.Fprintln(writer, dp[1])
	}
}
