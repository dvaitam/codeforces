package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type Edge struct{ to, id int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		var c int64
		fmt.Fscan(in, &n, &m, &c)

		adj := make([][]Edge, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], Edge{v, i})
			adj[v] = append(adj[v], Edge{u, i})
		}

		timer := 0
		tin := make([]int, n)
		low := make([]int, n)
		vis := make([]bool, n)
		compID := make([]int, n)
		var compSizes []int
		var bridgeSub [][]int

		var dfs func(int, int, int) int
		dfs = func(u, pe, cid int) int {
			vis[u] = true
			compID[u] = cid
			timer++
			tin[u] = timer
			low[u] = timer
			size := 1
			for _, e := range adj[u] {
				if e.id == pe {
					continue
				}
				v := e.to
				if !vis[v] {
					sub := dfs(v, e.id, cid)
					if low[v] > tin[u] {
						bridgeSub[cid] = append(bridgeSub[cid], sub)
					}
					if low[v] < low[u] {
						low[u] = low[v]
					}
					size += sub
				} else {
					if tin[v] < low[u] {
						low[u] = tin[v]
					}
				}
			}
			return size
		}

		for i := 0; i < n; i++ {
			if !vis[i] {
				cid := len(compSizes)
				bridgeSub = append(bridgeSub, []int{})
				size := dfs(i, -1, cid)
				compSizes = append(compSizes, size)
			}
		}

		k := len(compSizes)
		if k == 1 && len(bridgeSub[0]) == 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		// frequency of component sizes
		freq := make(map[int]int)
		for _, s := range compSizes {
			freq[s]++
		}

		maxN := n + 1
		// helper to shift in place
		addShift := func(dp []uint64, shift int) {
			word := shift / 64
			bit := shift % 64
			if word >= len(dp) {
				return
			}
			if bit == 0 {
				for i := len(dp) - 1; i >= word; i-- {
					dp[i] |= dp[i-word]
				}
			} else {
				for i := len(dp) - 1; i > word; i-- {
					dp[i] |= dp[i-word]<<bit | dp[i-word-1]>>(64-bit)
				}
				dp[word] |= dp[0] << bit
			}
		}

		computeDP := func(f map[int]int) []uint64 {
			dp := make([]uint64, (maxN+63)/64)
			dp[0] = 1
			for sz, cnt := range f {
				if cnt == 0 {
					continue
				}
				k := 1
				for cnt > 0 {
					take := k
					if take > cnt {
						take = cnt
					}
					shift := sz * take
					addShift(dp, shift)
					cnt -= take
					k <<= 1
				}
			}
			return dp
		}

		dpAll := computeDP(freq)

		// collect sizes that need dpWithout
		need := make(map[int]bool)
		for i, lst := range bridgeSub {
			if len(lst) > 0 {
				need[compSizes[i]] = true
			}
		}

		dpWithout := make(map[int][]uint64)
		for sz := range need {
			freq[sz]--
			dpWithout[sz] = computeDP(freq)
			freq[sz]++
		}

		bitsetTest := func(bs []uint64, pos int) bool {
			if pos < 0 || pos/64 >= len(bs) {
				return false
			}
			return (bs[pos/64]>>uint(pos%64))&1 == 1
		}

		searchClosest := func(bs []uint64, target, max int) int {
			if target < 0 {
				target = 0
			}
			if target > max {
				target = max
			}
			l := len(bs)
			word := target / 64
			bit := target % 64
			best := -1
			bestDist := max + 1
			if word < l {
				w := bs[word] & (^uint64(0) << uint(bit))
				if w != 0 {
					idx := bits.TrailingZeros64(w)
					cand := word*64 + idx
					if cand <= max {
						best = cand
						bestDist = cand - target
					}
				}
				w = bs[word] & ((uint64(1) << uint(bit)) - 1)
				if w != 0 {
					idx := 63 - bits.LeadingZeros64(w)
					cand := word*64 + idx
					d := target - cand
					if d < 0 {
						d = -d
					}
					if d < bestDist {
						bestDist = d
						best = cand
					}
				}
			}
			for i := word + 1; i < l && bestDist > 0; i++ {
				if bs[i] != 0 {
					idx := bits.TrailingZeros64(bs[i])
					cand := i*64 + idx
					if cand > max {
						cand = max
					}
					d := cand - target
					if d < 0 {
						d = -d
					}
					if d < bestDist {
						bestDist = d
						best = cand
					}
					break
				}
			}
			for i := word - 1; i >= 0 && bestDist > 0; i-- {
				if bs[i] != 0 {
					idx := 63 - bits.LeadingZeros64(bs[i])
					cand := i*64 + idx
					if cand > max {
						cand = max
					}
					d := target - cand
					if d < 0 {
						d = -d
					}
					if d < bestDist {
						bestDist = d
						best = cand
					}
					break
				}
			}
			if best == -1 {
				return 0
			}
			return best
		}

		bestXY := -1
		if k >= 2 {
			for s := 1; s < n; s++ {
				if bitsetTest(dpAll, s) {
					val := s * (n - s)
					if val > bestXY {
						bestXY = val
					}
				}
			}
		}

		for cid, lst := range bridgeSub {
			if len(lst) == 0 {
				continue
			}
			sz := compSizes[cid]
			dp := dpWithout[sz]
			maxS := n - sz
			for _, sub := range lst {
				parts := []int{sub, sz - sub}
				for _, p := range parts {
					if p <= 0 || p >= n {
						continue
					}
					target := n/2 - p
					s := searchClosest(dp, target, maxS)
					x := p + s
					if x <= 0 || x >= n {
						continue
					}
					val := x * (n - x)
					if val > bestXY {
						bestXY = val
					}
				}
			}
		}

		if bestXY < 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		cost := int64(n*n) - 2*int64(bestXY) + int64(k-1)*c
		fmt.Fprintln(out, cost)
	}
}
