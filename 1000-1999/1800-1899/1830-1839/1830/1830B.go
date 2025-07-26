package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		freq := make(map[int]map[int]int)
		seen := make(map[int]bool)
		uniqueA := []int{}
		for i := 0; i < n; i++ {
			ai := a[i]
			bi := b[i]
			if freq[ai] == nil {
				freq[ai] = make(map[int]int)
			}
			freq[ai][bi]++
			if !seen[ai] {
				seen[ai] = true
				uniqueA = append(uniqueA, ai)
			}
		}
		sort.Ints(uniqueA)
		limit := 2 * n
		var ans int64
		for ix := 0; ix < len(uniqueA); ix++ {
			x := uniqueA[ix]
			mapX := freq[x]
			for iy := ix; iy < len(uniqueA); iy++ {
				y := uniqueA[iy]
				if x*y > limit {
					break
				}
				mapY := freq[y]
				p := x * y
				if x == y {
					for b1, c1 := range mapX {
						b2 := p - b1
						if b2 < 1 || b2 > n {
							continue
						}
						if b2 < b1 {
							continue
						}
						if c2, ok := mapX[b2]; ok {
							if b1 == b2 {
								ans += int64(c1) * int64(c1-1) / 2
							} else {
								ans += int64(c1) * int64(c2)
							}
						}
					}
				} else {
					if len(mapX) <= len(mapY) {
						for b1, c1 := range mapX {
							b2 := p - b1
							if b2 < 1 || b2 > n {
								continue
							}
							if c2 := mapY[b2]; c2 > 0 {
								ans += int64(c1) * int64(c2)
							}
						}
					} else {
						for b2, c2 := range mapY {
							b1 := p - b2
							if b1 < 1 || b1 > n {
								continue
							}
							if c1 := mapX[b1]; c1 > 0 {
								ans += int64(c1) * int64(c2)
							}
						}
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
