package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cycle struct {
	mask  int
	nodes []int64
}

type Move struct {
	num  int64
	dest int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	boxes := make([][]int64, k)
	boxsum := make([]int64, k)
	whichbox := make(map[int64]int)
	var total int64
	var nodes []int64

	for i := 0; i < k; i++ {
		var n int
		fmt.Fscan(in, &n)
		boxes[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &boxes[i][j])
			v := boxes[i][j]
			whichbox[v] = i
			nodes = append(nodes, v)
			total += v
			boxsum[i] += v
		}
	}
	if total%int64(k) != 0 {
		fmt.Fprintln(out, "NO")
		return
	}
	target := total / int64(k)
	nxt := make(map[int64]int64, len(nodes))
	for i := 0; i < k; i++ {
		for _, v := range boxes[i] {
			needed := target - boxsum[i] + v
			if j, ok := whichbox[needed]; ok {
				nxt[v] = needed
			} else {
				nxt[v] = -1
			}
		}
	}

	processed := make(map[int64]bool)
	var validcycles []Cycle

	for _, start := range nodes {
		if processed[start] {
			continue
		}
		position := make(map[int64]int)
		var path []int64
		cur := start
		found := false
		for cur != -1 {
			if _, seen := position[cur]; seen {
				found = true
				break
			}
			position[cur] = len(path)
			path = append(path, cur)
			cur = nxt[cur]
		}
		if found {
			pos := position[cur]
			var mask int
			var cycle []int64
			ok := true
			for _, v := range path[pos:] {
				b := whichbox[v]
				if mask&(1<<b) != 0 {
					ok = false
					break
				}
				mask |= 1 << b
				cycle = append(cycle, v)
			}
			if ok {
				validcycles = append(validcycles, Cycle{mask: mask, nodes: cycle})
			}
		}
		for _, v := range path {
			processed[v] = true
		}
	}

	fullMask := (1 << k) - 1
	// group cycles by mask
	cyclesByMask := make([][]int, fullMask+1)
	for i, c := range validcycles {
		cyclesByMask[c.mask] = append(cyclesByMask[c.mask], i)
	}
	// dp over subsets
	dp := make([]bool, fullMask+1)
	parent := make([]int, fullMask+1)
	used := make([]int, fullMask+1)
	f := make([]bool, fullMask+1)
	dp[0] = true
	for _, c := range validcycles {
		f[c.mask] = true
	}
	for mask := 0; mask <= fullMask; mask++ {
		if !dp[mask] {
			continue
		}
		remain := fullMask ^ mask
		for s := remain; s > 0; s = (s - 1) & remain {
			if !f[s] {
				continue
			}
			newMask := mask | s
			if !dp[newMask] {
				dp[newMask] = true
				parent[newMask] = mask
				used[newMask] = s
			}
		}
	}
	if !dp[fullMask] {
		fmt.Fprintln(out, "NO")
		return
	}
	fmt.Fprintln(out, "YES")
	// reconstruct
	var chain []int
	for m := fullMask; m != 0; m = parent[m] {
		chain = append(chain, used[m])
	}
	// reverse chain
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	var ans []Move
	for _, sm := range chain {
		ci := cyclesByMask[sm][0]
		cycle := validcycles[ci].nodes
		// reverse to match order
		for i, j := 0, len(cycle)-1; i < j; i, j = i+1, j-1 {
			cycle[i], cycle[j] = cycle[j], cycle[i]
		}
		t := len(cycle)
		for i := 0; i < t; i++ {
			num := cycle[i]
			destBox := whichbox[cycle[(i+1)%t]] + 1
			ans = append(ans, Move{num: num, dest: destBox})
		}
	}
	// sort by source box index
	sort.Slice(ans, func(i, j int) bool {
		bi := whichbox[ans[i].num]
		bj := whichbox[ans[j].num]
		return bi < bj
	})
	for _, mv := range ans {
		fmt.Fprintf(out, "%d %d\n", mv.num, mv.dest)
	}
}
