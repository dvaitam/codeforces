package main

import (
	"bufio"
	"fmt"
	"os"
)

func minAdditional(nodes int, cap int, depth int) int64 {
	var sum int64
	rem := nodes
	c := cap
	d := depth
	for rem > 0 {
		use := c
		if use > rem {
			use = rem
		}
		sum += int64(use * d)
		rem -= use
		c = use * 2
		d++
	}
	return sum
}

func maxAdditional(nodes int, depth int) int64 {
	n := int64(nodes)
	return n*int64(depth) + n*(n-1)/2
}

func constructTree(n, d int) ([]int, bool) {
	maxSum := int64(n*(n-1)) / 2
	minSum := minAdditional(n-1, 2, 1)
	if int64(d) < minSum || int64(d) > maxSum {
		return nil, false
	}

	levels := []int{1}
	curSum := int64(0)
	rem := n - 1
	avail := 1
	depth := 1
	for rem > 0 {
		maxNodes := avail * 2
		if maxNodes > rem {
			maxNodes = rem
		}
		chosen := 0
		for x := 1; x <= maxNodes; x++ {
			remaining := rem - x
			minAdd := minAdditional(remaining, x*2, depth+1)
			maxAdd := maxAdditional(remaining, depth+1)
			totalMin := curSum + int64(x*depth) + minAdd
			totalMax := curSum + int64(x*depth) + maxAdd
			if int64(d) >= totalMin && int64(d) <= totalMax {
				chosen = x
				break
			}
		}
		if chosen == 0 {
			return nil, false
		}
		levels = append(levels, chosen)
		curSum += int64(chosen * depth)
		rem -= chosen
		avail = chosen
		depth++
	}

	parent := make([]int, n+1)
	parent[1] = 0
	prev := []int{1}
	idx := 2
	for lvl := 1; lvl < len(levels); lvl++ {
		cnt := levels[lvl]
		next := make([]int, 0, cnt)
		pIdx := 0
		used := 0
		for i := 0; i < cnt; i++ {
			if used == 2 {
				pIdx++
				used = 0
			}
			if pIdx >= len(prev) || idx > n {
				return nil, false
			}
			parent[idx] = prev[pIdx]
			next = append(next, idx)
			idx++
			used++
		}
		prev = next
	}

	return parent, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, d int
		fmt.Fscan(in, &n, &d)
		parent, ok := constructTree(n, d)
		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}
		fmt.Fprintln(out, "YES")
		for i := 2; i <= n; i++ {
			if i > 2 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, parent[i])
		}
		if n >= 2 {
			fmt.Fprintln(out)
		}
	}
}
