package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	d int
	t int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		a := make([]int, n)
		b := make([]int, n)
		s := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i], &s[i])
		}

		city := make([]int, n)
		for cityID := 1; cityID <= m; cityID++ {
			var k int
			fmt.Fscan(in, &k)
			for j := 0; j < k; j++ {
				var idx int
				fmt.Fscan(in, &idx)
				city[idx-1] = cityID
			}
		}

		maxStrengthCity2 := -1
		bMaxCity2 := make(map[int]int)
		minBCity1 := make(map[int]int)
		specOrder := make([]int, 0)

		for i := 0; i < n; i++ {
			if city[i] == 1 {
				if _, ok := minBCity1[s[i]]; !ok {
					minBCity1[s[i]] = b[i]
					specOrder = append(specOrder, s[i])
				} else if b[i] < minBCity1[s[i]] {
					minBCity1[s[i]] = b[i]
				}
			} else {
				if a[i] > maxStrengthCity2 {
					maxStrengthCity2 = a[i]
				}
				if val, ok := bMaxCity2[s[i]]; !ok || b[i] > val {
					bMaxCity2[s[i]] = b[i]
				}
			}
		}

		possible := true
		solutions := make([]pair, 0, len(specOrder))
		for _, spec := range specOrder {
			limit := max(maxStrengthCity2, bMaxCity2[spec])
			if minBCity1[spec] <= limit {
				possible = false
				break
			}
			difficulty := limit + 1
			solutions = append(solutions, pair{difficulty, spec})
		}

		if !possible {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(solutions))
		for _, sol := range solutions {
			fmt.Fprintf(out, "%d %d\n", sol.d, sol.t)
		}
	}
}
