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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	const V = 200
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		// positions for each value
		pos := make([][]int, V+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			v := arr[i]
			pos[v] = append(pos[v], i+1) // 1-indexed
		}
		// prefix counts
		pref := make([][]int, V+1)
		for v := 1; v <= V; v++ {
			pref[v] = make([]int, n+1)
		}
		for i := 1; i <= n; i++ {
			val := arr[i-1]
			for v := 1; v <= V; v++ {
				pref[v][i] = pref[v][i-1]
			}
			pref[val][i]++
		}
		ans := 0
		// case when l=0 -> just middle block
		for v := 1; v <= V; v++ {
			if len(pos[v]) > ans {
				ans = len(pos[v])
			}
		}
		// consider blocks with value x at the ends
		for x := 1; x <= V; x++ {
			px := pos[x]
			m := len(px)
			for l := 1; l*2 <= m; l++ {
				left := px[l-1]
				right := px[m-l]
				maxMid := 0
				for y := 1; y <= V; y++ {
					cnt := pref[y][right-1] - pref[y][left]
					if cnt > maxMid {
						maxMid = cnt
					}
				}
				tot := 2*l + maxMid
				if tot > ans {
					ans = tot
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
