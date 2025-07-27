package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	block int
	coord int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		X := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &X[i])
		}
		Y := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &Y[i])
		}
		setX := make(map[int]struct{}, n)
		setY := make(map[int]struct{}, m)
		for _, v := range X {
			setX[v] = struct{}{}
		}
		for _, v := range Y {
			setY[v] = struct{}{}
		}

		vertical := make([]pair, 0, k)
		horizontal := make([]pair, 0, k)
		for i := 0; i < k; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			_, inX := setX[x]
			_, inY := setY[y]
			if inX && inY {
				continue
			}
			if inX {
				block := sort.Search(len(Y), func(j int) bool { return Y[j] > y })
				vertical = append(vertical, pair{block, x})
			} else {
				block := sort.Search(len(X), func(j int) bool { return X[j] > x })
				horizontal = append(horizontal, pair{block, y})
			}
		}

		ans := int64(0)
		sort.Slice(vertical, func(i, j int) bool {
			if vertical[i].block == vertical[j].block {
				return vertical[i].coord < vertical[j].coord
			}
			return vertical[i].block < vertical[j].block
		})
		for i := 0; i < len(vertical); {
			j := i
			for j < len(vertical) && vertical[j].block == vertical[i].block {
				j++
			}
			cnt := int64(j - i)
			ans += cnt * (cnt - 1) / 2
			for l := i; l < j; {
				r := l
				for r < j && vertical[r].coord == vertical[l].coord {
					r++
				}
				c := int64(r - l)
				ans -= c * (c - 1) / 2
				l = r
			}
			i = j
		}

		sort.Slice(horizontal, func(i, j int) bool {
			if horizontal[i].block == horizontal[j].block {
				return horizontal[i].coord < horizontal[j].coord
			}
			return horizontal[i].block < horizontal[j].block
		})
		for i := 0; i < len(horizontal); {
			j := i
			for j < len(horizontal) && horizontal[j].block == horizontal[i].block {
				j++
			}
			cnt := int64(j - i)
			ans += cnt * (cnt - 1) / 2
			for l := i; l < j; {
				r := l
				for r < j && horizontal[r].coord == horizontal[l].coord {
					r++
				}
				c := int64(r - l)
				ans -= c * (c - 1) / 2
				l = r
			}
			i = j
		}

		fmt.Fprintln(out, ans)
	}
}
