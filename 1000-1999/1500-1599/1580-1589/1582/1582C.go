package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

func cost(s string, ch byte) int {
	l := 0
	r := len(s) - 1
	cnt := 0
	for l < r {
		if s[l] == s[r] {
			l++
			r--
		} else if s[l] == ch {
			cnt++
			l++
		} else if s[r] == ch {
			cnt++
			r--
		} else {
			return INF
		}
	}
	return cnt
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		best := INF
		for c := byte('a'); c <= byte('z'); c++ {
			v := cost(s, c)
			if v < best {
				best = v
			}
		}
		if best == INF {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, best)
		}
	}
}
