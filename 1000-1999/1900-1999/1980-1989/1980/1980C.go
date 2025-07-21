package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var m int
		fmt.Fscan(in, &m)
		cnt := make(map[int]int)
		for i := 0; i < m; i++ {
			var d int
			fmt.Fscan(in, &d)
			cnt[d]++
		}
		need := make(map[int]int)
		possible := true
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				need[b[i]]++
			}
		}
		for val, c := range need {
			if cnt[val] < c {
				possible = false
				break
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
