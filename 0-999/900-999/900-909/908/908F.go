package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int64, n)
	c := make([]byte, n)
	for i := 0; i < n; i++ {
		var color string
		fmt.Fscan(in, &p[i], &color)
		c[i] = color[0]
	}
	greens := []int{}
	for i := 0; i < n; i++ {
		if c[i] == 'G' {
			greens = append(greens, i)
		}
	}
	var ans int64 = 0
	if len(greens) == 0 {
		var firstR, lastR, firstB, lastB int64
		firstR, lastR, firstB, lastB = -1, -1, -1, -1
		for i := 0; i < n; i++ {
			if c[i] == 'R' {
				if firstR == -1 {
					firstR = p[i]
				}
				lastR = p[i]
			} else if c[i] == 'B' {
				if firstB == -1 {
					firstB = p[i]
				}
				lastB = p[i]
			}
		}
		if firstR != -1 {
			ans += lastR - firstR
		}
		if firstB != -1 {
			ans += lastB - firstB
		}
		if firstR != -1 && firstB != -1 {
			minDiff := int64(1<<63 - 1)
			for i := 1; i < n; i++ {
				if (c[i] == 'R' && c[i-1] == 'B') || (c[i] == 'B' && c[i-1] == 'R') {
					diff := p[i] - p[i-1]
					if diff < minDiff {
						minDiff = diff
					}
				}
			}
			ans += minDiff
		}
		fmt.Println(ans)
		return
	}

	firstG := greens[0]
	lastG := greens[len(greens)-1]
	var firstR int64 = -1
	var firstB int64 = -1
	for i := 0; i < firstG; i++ {
		if c[i] == 'R' {
			if firstR == -1 {
				firstR = p[i]
			}
		} else if c[i] == 'B' {
			if firstB == -1 {
				firstB = p[i]
			}
		}
	}
	if firstR != -1 {
		ans += p[firstG] - firstR
	}
	if firstB != -1 {
		ans += p[firstG] - firstB
	}
	var lastR int64 = -1
	var lastB int64 = -1
	for i := n - 1; i > lastG; i-- {
		if c[i] == 'R' {
			if lastR == -1 {
				lastR = p[i]
			}
		} else if c[i] == 'B' {
			if lastB == -1 {
				lastB = p[i]
			}
		}
	}
	if lastR != -1 {
		ans += lastR - p[lastG]
	}
	if lastB != -1 {
		ans += lastB - p[lastG]
	}

	for idx := 0; idx < len(greens)-1; idx++ {
		g1 := greens[idx]
		g2 := greens[idx+1]
		attach := int64(0)
		for j := g1 + 1; j < g2; j++ {
			d1 := p[j] - p[g1]
			d2 := p[g2] - p[j]
			if d1 < d2 {
				attach += d1
			} else {
				attach += d2
			}
		}
		d := p[g2] - p[g1]
		cost1 := d + attach
		cost2 := 2 * d
		if cost1 < cost2 {
			ans += cost1
		} else {
			ans += cost2
		}
	}

	fmt.Println(ans)
}
