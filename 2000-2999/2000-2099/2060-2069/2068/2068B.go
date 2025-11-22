package main

import (
	"bufio"
	"fmt"
	"os"
)

func tri(x int64) int64 {
	return x * (x - 1) / 2
}

type pair struct {
	h, w int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int64
	fmt.Fscan(in, &k)

	if k == 0 {
		fmt.Fprintln(out, "1 1")
		fmt.Fprintln(out, ".")
		return
	}

	const limit = 2025

	// Precompute all possible values of comb(h,2)*comb(w,2) with h,w >=2
	valMap := make(map[int64]pair, 4_500_000)
	for h := 2; h <= limit; h++ {
		th := tri(int64(h))
		for w := 2; w <= limit; w++ {
			v := th * tri(int64(w))
			if _, ok := valMap[v]; !ok {
				valMap[v] = pair{h, w}
			}
		}
	}

	// Try single block
	if p, ok := valMap[k]; ok {
		fmt.Fprintf(out, "%d %d\n", p.h, p.w)
		for i := 0; i < p.h; i++ {
			for j := 0; j < p.w; j++ {
				fmt.Fprint(out, "#")
			}
			fmt.Fprintln(out)
		}
		return
	}

	// Try two blocks stacked vertically with one separator row of dots.
	type block struct {
		h, w int
	}
	blocks := []block{}
	found := false
	var b1, b2 block

	for h1 := 2; h1 <= limit && !found; h1++ {
		th1 := tri(int64(h1))
		for w1 := 2; w1 <= limit; w1++ {
			v1 := th1 * tri(int64(w1))
			if v1 > k {
				break
			}
			remaining := k - v1
			if remaining == 0 {
				b1 = block{h1, w1}
				found = true
				break
			}
			if p2, ok := valMap[remaining]; ok {
				if h1+p2.h+1 <= limit {
					b1 = block{h1, w1}
					b2 = block{p2.h, p2.w}
					found = true
					break
				}
			}
		}
	}

	if !found {
		// Greedy fallback: build blocks top-down while height budget allows.
		heightLeft := limit
		rem := k
		for rem > 0 && heightLeft >= 2 {
			// If not the first block, reserve a separator row.
			if len(blocks) > 0 {
				heightLeft--
				if heightLeft < 2 {
					break
				}
			}
			// Choose the largest height whose maximum possible value does not exceed rem.
			h := 0
			for cand := heightLeft; cand >= 2; cand-- {
				if tri(int64(cand))*tri(limit) <= rem {
					h = cand
					break
				}
			}
			if h == 0 {
				h = heightLeft
			}
			th := tri(int64(h))
			// Choose maximal width with value not exceeding remaining.
			w := 2
			for w <= limit {
				v := th * tri(int64(w))
				if v > rem {
					break
				}
				w++
			}
			w--
			if w < 2 {
				break
			}
			val := th * tri(int64(w))
			blocks = append(blocks, block{h, w})
			rem -= val
			heightLeft -= h
		}
		found = rem == 0
		if !found {
			return
		}
	} else {
		blocks = []block{b1}
		if b2.h > 0 {
			blocks = append(blocks, b2)
		}
	}

	// Compute final dimensions.
	totalH := 0
	maxW := 0
	for i, b := range blocks {
		totalH += b.h
		if i+1 < len(blocks) {
			totalH++ // separator row
		}
		if b.w > maxW {
			maxW = b.w
		}
	}
	if totalH == 0 || maxW == 0 || totalH > limit || maxW > limit {
		// Should not happen with provided constraints.
		return
	}

	fmt.Fprintf(out, "%d %d\n", totalH, maxW)
	for idx, b := range blocks {
		for i := 0; i < b.h; i++ {
			for j := 0; j < maxW; j++ {
				if j < b.w {
					fmt.Fprint(out, "#")
				} else {
					fmt.Fprint(out, ".")
				}
			}
			fmt.Fprintln(out)
		}
		if idx+1 < len(blocks) {
			for j := 0; j < maxW; j++ {
				fmt.Fprint(out, ".")
			}
			fmt.Fprintln(out)
		}
	}
}
