package main

import (
	"bufio"
	"fmt"
	"os"
)

func sortPair(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	triples := make([][3]int, n-2)
	cnt := make([]int, n+1)
	pairMap := make(map[[2]int][]int)
	contain := make([][]int, n+1)

	for i := 0; i < n-2; i++ {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		triples[i] = [3]int{a, b, c}
		cnt[a]++
		cnt[b]++
		cnt[c]++
		contain[a] = append(contain[a], i)
		contain[b] = append(contain[b], i)
		contain[c] = append(contain[c], i)

		arr := [][3]int{{a, b, c}, {a, c, b}, {b, c, a}}
		for _, p := range arr {
			u, v, w := p[0], p[1], p[2]
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			pairMap[key] = append(pairMap[key], w)
		}
	}

	start := 0
	second := 0
	third := 0
	for i := 1; i <= n; i++ {
		if cnt[i] == 1 {
			idx := contain[i][0]
			t := triples[idx]
			for j := 0; j < 3; j++ {
				x := t[j]
				if x != i && cnt[x] == 2 {
					start = i
					second = x
					for k := 0; k < 3; k++ {
						y := t[k]
						if y != i && y != x {
							third = y
							break
						}
					}
					break
				}
			}
			if start != 0 {
				break
			}
		}
	}

	ans := make([]int, n)
	ans[0] = start
	ans[1] = second
	ans[2] = third

	for i := 3; i < n; i++ {
		u := ans[i-2]
		v := ans[i-1]
		if u > v {
			u, v = v, u
		}
		cand := pairMap[[2]int{u, v}]
		if len(cand) == 1 {
			ans[i] = cand[0]
		} else {
			if cand[0] == ans[i-3] {
				ans[i] = cand[1]
			} else {
				ans[i] = cand[0]
			}
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
