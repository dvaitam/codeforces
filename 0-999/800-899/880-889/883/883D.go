package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func canCover(t int, stars, packs []int) bool {
	m, n := len(stars), len(packs)
	if m == 0 {
		return true
	}
	i := 0
	for j := 0; j < n && i < m; j++ {
		p := packs[j]
		// if current star is left of this packman
		if stars[i] < p {
			// see if next packman can reach it
			if j+1 < n && packs[j+1]-stars[i] <= t {
				// orient this packman right to help with right side
				for i < m && stars[i] >= p && stars[i]-p <= t {
					i++
				}
			} else {
				// must orient left
				if p-stars[i] > t {
					return false
				}
				for i < m && stars[i] <= p && p-stars[i] <= t {
					i++
				}
			}
		} else {
			// star is to the right
			for i < m && stars[i] >= p && stars[i]-p <= t {
				i++
			}
		}
	}
	return i == m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	stars := make([]int, 0)
	packs := make([]int, 0)
	for i, ch := range s {
		if ch == '*' {
			stars = append(stars, i)
		} else if ch == 'P' {
			packs = append(packs, i)
		}
	}
	sort.Ints(stars)
	sort.Ints(packs)
	m := len(stars)
	k := len(packs)
	if k == 0 || m == 0 {
		fmt.Println(0, 0)
		return
	}
	if k == 1 {
		p := packs[0]
		// count stars on each side
		idx := sort.SearchInts(stars, p)
		left := idx
		right := m - idx
		if left >= right {
			time := 0
			if left > 0 {
				time = p - stars[0]
			}
			fmt.Println(left, time)
		} else {
			time := 0
			if right > 0 {
				time = stars[m-1] - p
			}
			fmt.Println(right, time)
		}
		return
	}
	// k >= 2, we can eat all stars
	l, r := 0, n
	for l < r {
		mid := (l + r) / 2
		if canCover(mid, stars, packs) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Println(m, l)
}
