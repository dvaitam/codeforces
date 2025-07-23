package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canEat(t int, stars, packs []int) bool {
	j := 0
	for _, p := range packs {
		if j >= len(stars) {
			return true
		}
		if stars[j] > p {
			// all remaining stars to the right
			reach := p + t
			for j < len(stars) && stars[j] <= reach {
				j++
			}
		} else {
			// first star is on or left of packman
			d := p - stars[j]
			if d > t {
				return false
			}
			reach := max(p+(t-2*d), p+(t-d)/2)
			for j < len(stars) && stars[j] <= reach {
				j++
			}
		}
	}
	return j >= len(stars)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var field string
	fmt.Fscan(in, &field)

	stars := make([]int, 0)
	packs := make([]int, 0)
	for i := 0; i < n; i++ {
		switch field[i] {
		case '*':
			stars = append(stars, i)
		case 'P':
			packs = append(packs, i)
		}
	}

	l, r := 0, 2*n
	for l < r {
		mid := (l + r) / 2
		if canEat(mid, stars, packs) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Println(l)
}
