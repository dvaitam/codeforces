package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ a, b int }

func check(H, W int, rects [][2]int) bool {
	area := 0
	for _, r := range rects {
		area += r[0] * r[1]
	}
	if area != H*W {
		return false
	}
	counts := make(map[pair]int, len(rects))
	byH := make(map[int]map[int]int)
	byW := make(map[int]map[int]int)
	for _, r := range rects {
		p := pair{r[0], r[1]}
		counts[p]++
		if _, ok := byH[r[0]]; !ok {
			byH[r[0]] = map[int]int{}
		}
		byH[r[0]][r[1]]++
		if _, ok := byW[r[1]]; !ok {
			byW[r[1]] = map[int]int{}
		}
		byW[r[1]][r[0]]++
	}
	remove := func(h, w int) {
		p := pair{h, w}
		c := counts[p] - 1
		if c == 0 {
			delete(counts, p)
			delete(byH[h], w)
			if len(byH[h]) == 0 {
				delete(byH, h)
			}
			delete(byW[w], h)
			if len(byW[w]) == 0 {
				delete(byW, w)
			}
		} else {
			counts[p] = c
			byH[h][w] = c
			byW[w][h] = c
		}
	}
	h := H
	w := W
	n := len(rects)
	for n > 0 {
		if m, ok := byH[h]; ok {
			var w1 int
			for w1 = range m {
				break
			}
			if w1 > w {
				return false
			}
			remove(h, w1)
			w -= w1
		} else if m, ok := byW[w]; ok {
			var h1 int
			for h1 = range m {
				break
			}
			if h1 > h {
				return false
			}
			remove(h1, w)
			h -= h1
		} else {
			return false
		}
		n--
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		rects := make([][2]int, n)
		var sum int
		maxH, maxW := 0, 0
		for i := 0; i < n; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			rects[i] = [2]int{a, b}
			sum += a * b
			if a > maxH {
				maxH = a
			}
			if b > maxW {
				maxW = b
			}
		}
		ans := make([][2]int, 0, 2)
		if sum%maxH == 0 {
			W := sum / maxH
			if check(maxH, W, rects) {
				ans = append(ans, [2]int{maxH, W})
			}
		}
		if sum%maxW == 0 {
			H := sum / maxW
			if check(H, maxW, rects) {
				pair := [2]int{H, maxW}
				if len(ans) == 0 || ans[0] != pair {
					ans = append(ans, pair)
				}
			}
		}
		fmt.Fprintln(out, len(ans))
		for _, p := range ans {
			fmt.Fprintln(out, p[0], p[1])
		}
	}
}
