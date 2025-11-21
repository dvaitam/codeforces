package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	r int
	c int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, n, m int
		fmt.Fscan(in, &l, &n, &m)

		a := make([]int, l)
		for i := 0; i < l; i++ {
			fmt.Fscan(in, &a[i])
		}

		pos := make(map[int][]pair)
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				var val int
				fmt.Fscan(in, &val)
				pos[val] = append(pos[val], pair{r: i, c: j})
			}
		}

		const inf = int(1e9)
		curRowMin := make([]int, n+2)
		for i := range curRowMin {
			curRowMin[i] = inf
		}

		length := 0

		for idx := 0; idx < l; idx++ {
			points := pos[a[idx]]
			if len(points) == 0 {
				break
			}

			if length == 0 {
				for _, p := range points {
					if p.c < curRowMin[p.r] {
						curRowMin[p.r] = p.c
					}
				}
				possible := false
				for r := 1; r <= n; r++ {
					if curRowMin[r] != inf {
						possible = true
						break
					}
				}
				if !possible {
					break
				}
				length = 1
				continue
			}

			pref := make([]int, n+2)
			pref[0] = inf
			for r := 1; r <= n; r++ {
				pref[r] = pref[r-1]
				if curRowMin[r] < pref[r] {
					pref[r] = curRowMin[r]
				}
			}

			newRowMin := make([]int, n+2)
			for r := range newRowMin {
				newRowMin[r] = inf
			}
			any := false

			for _, p := range points {
				if pref[p.r-1] < p.c {
					if p.c < newRowMin[p.r] {
						newRowMin[p.r] = p.c
					}
					any = true
				}
			}

			if !any {
				break
			}

			curRowMin = newRowMin
			length++
		}

		if length%2 == 1 {
			fmt.Fprintln(out, "T")
		} else {
			fmt.Fprintln(out, "N")
		}
	}
}
