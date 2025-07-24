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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		bits := make([][]int, n)
		const M = 200000 + 1
		cnt := make([]int, M)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(in, &k)
			bits[i] = make([]int, k)
			for j := 0; j < k; j++ {
				var p int
				fmt.Fscan(in, &p)
				bits[i][j] = p
				cnt[p]++
			}
		}
		ans := "No"
		for i := 0; i < n && ans == "No"; i++ {
			ok := true
			for _, b := range bits[i] {
				if cnt[b] < 2 {
					ok = false
					break
				}
			}
			if ok {
				ans = "Yes"
			}
		}
		fmt.Fprintln(out, ans)
	}
}
