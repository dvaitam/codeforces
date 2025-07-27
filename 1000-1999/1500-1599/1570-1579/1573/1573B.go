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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		posA := make([]int, 2*n+3)
		posB := make([]int, 2*n+3)
		for i, v := range a {
			posA[v] = i + 1
		}
		for i, v := range b {
			posB[v] = i + 1
		}
		const INF = int(1e9)
		bestEven := make([]int, 2*n+5)
		for e := 2 * n; e >= 2; e -= 2 {
			cur := posB[e]
			if cur == 0 {
				cur = INF
			}
			next := bestEven[e+2]
			if next == 0 {
				next = INF
			}
			if cur < next {
				bestEven[e] = cur
			} else {
				bestEven[e] = next
			}
		}
		ans := INF
		for o := 1; o < 2*n; o += 2 {
			pa := posA[o]
			if pa == 0 {
				pa = INF
			}
			pb := bestEven[o+1]
			if pa+pb < ans {
				ans = pa + pb
			}
		}
		fmt.Fprintln(writer, ans-2)
	}
}
