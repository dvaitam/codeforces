package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	strengths := make([]int, 3)
	fmt.Fscan(in, &strengths[0], &strengths[1], &strengths[2])
	sort.Ints(strengths)
	a, b, c := strengths[0], strengths[1], strengths[2]

	criminals := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &criminals[i])
	}
	sort.Ints(criminals)

	used := make([]bool, n)
	remain := n

	pa := n - 1
	for pa >= 0 && criminals[pa] > a {
		pa--
	}
	pb := n - 1
	for pb >= 0 && criminals[pb] > b {
		pb--
	}
	r := n - 1
	hours := 0

	remove := func(limit int, p *int) bool {
		for *p >= 0 && (used[*p] || criminals[*p] > limit) {
			(*p)--
		}
		if *p >= 0 {
			used[*p] = true
			remain--
			(*p)--
			return true
		}
		return false
	}

	for remain > 0 {
		for r >= 0 && used[r] {
			r--
		}
		if r < 0 {
			break
		}
		idx := r
		x := criminals[idx]
		used[idx] = true
		remain--
		r--
		hours++

		if x > a+b+c {
			fmt.Println(-1)
			return
		}
		if x > b+c {
			continue
		}
		if x > c {
			if x > a+c {
				remove(a, &pa)
			} else {
				remove(b, &pb)
			}
			continue
		}
		remove(b, &pb)
		remove(a, &pa)
	}

	fmt.Println(hours)
}
