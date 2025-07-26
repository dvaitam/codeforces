package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 1000000

var prefix [maxN + 1]int64
var tri [2024]int

func init() {
	for i := 1; i <= maxN; i++ {
		prefix[i] = prefix[i-1] + int64(i*i)
	}
	for i := 1; i < len(tri); i++ {
		tri[i] = tri[i-1] + i
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		r := 1
		for tri[r] < n {
			r++
		}
		c := n - tri[r-1]
		var ans int64
		for i := 1; i <= r; i++ {
			j1 := c - (r - i)
			if j1 < 1 {
				j1 = 1
			}
			if j1 > i {
				continue
			}
			j2 := c
			if j2 > i {
				j2 = i
			}
			start := tri[i-1] + j1
			end := tri[i-1] + j2
			ans += prefix[end] - prefix[start-1]
		}
		fmt.Fprintln(out, ans)
	}
}
