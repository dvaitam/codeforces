package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	// adjacency for bfs and independent set counting
	adjList := make([][]int, n)

	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		adjList[x] = append(adjList[x], y)
		adjList[y] = append(adjList[y], x)
	}

	// count connected components, isolated vertices and check bipartite
	visited := make([]bool, n)
	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	cc := 0
	iso := 0
	bipartite := true

	for i := 0; i < n; i++ {
		if !visited[i] {
			cc++
			if len(adjList[i]) == 0 {
				iso++
				visited[i] = true
				color[i] = 0
				continue
			}
			// bfs/dfs
			stack := []int{i}
			visited[i] = true
			color[i] = 0
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, to := range adjList[v] {
					if !visited[to] {
						visited[to] = true
						color[to] = color[v] ^ 1
						stack = append(stack, to)
					} else if color[to] == color[v] {
						bipartite = false
					}
				}
			}
		}
	}

	// meet-in-the-middle for number of independent sets
	var n1 = n / 2
	var n2 = n - n1

	// adjacency split
	adj1 := make([]uint32, n1)
	adj2 := make([]uint32, n2)
	cross2 := make([]uint32, n2)

	for v := 0; v < n; v++ {
		if v < n1 {
			for _, u := range adjList[v] {
				if u < n1 {
					adj1[v] |= 1 << uint(u)
				} else {
					cross2[u-n1] |= 1 << uint(v)
				}
			}
		} else {
			v2 := v - n1
			for _, u := range adjList[v] {
				if u < n1 {
					cross2[v2] |= 1 << uint(u)
				} else {
					adj2[v2] |= 1 << uint(u-n1)
				}
			}
		}
	}

	// first half independent flags
	size1 := 1 << uint(n1)
	ind1 := make([]bool, size1)
	ind1[0] = true
	for mask := 1; mask < size1; mask++ {
		lb := mask & -mask
		i := 0
		for (lb>>uint(i))&1 == 0 {
			i++
		}
		prev := mask ^ lb
		if ind1[prev] && (int(adj1[i])&prev) == 0 {
			ind1[mask] = true
		}
	}

	// second half dynamic
	size2 := 1 << uint(n2)
	ind2 := make([]bool, size2)
	crossMask2 := make([]uint32, size2)
	ind2[0] = true
	// counts by cross mask to first half
	cnt := make([]int64, size1)
	cnt[0] = 1
	for mask := 1; mask < size2; mask++ {
		lb := mask & -mask
		i := 0
		for (lb>>uint(i))&1 == 0 {
			i++
		}
		prev := mask ^ lb
		if ind2[prev] && (int(adj2[i])&prev) == 0 {
			ind2[mask] = true
			crossMask2[mask] = crossMask2[prev] | cross2[i]
			cnt[crossMask2[mask]]++
		}
	}

	// subset sum DP on counts
	for i := 0; i < n1; i++ {
		for mask := 0; mask < size1; mask++ {
			if mask&(1<<uint(i)) != 0 {
				cnt[mask] += cnt[mask^(1<<uint(i))]
			}
		}
	}

	// count independent sets by combining
	var indepSets int64
	maskAll := size1 - 1
	for mask := 0; mask < size1; mask++ {
		if ind1[mask] {
			indepSets += cnt[maskAll^mask]
		}
	}

	// compute powers of two
	pow2 := func(x int) int64 {
		return int64(1) << uint(x)
	}

	total := pow2(n)
	countS0 := indepSets
	countS2 := indepSets
	countS1 := pow2(cc)

	countS0S1 := pow2(iso)
	countS1S2 := pow2(iso)
	var countS0S2 int64
	if bipartite {
		countS0S2 = pow2(cc)
	} else {
		countS0S2 = 0
	}

	var countAll int64
	if m == 0 {
		countAll = pow2(n)
	} else {
		countAll = 0
	}

	result := total - (countS0 + countS1 + countS2) + (countS0S1 + countS0S2 + countS1S2) - countAll
	fmt.Fprintln(out, result)
}
