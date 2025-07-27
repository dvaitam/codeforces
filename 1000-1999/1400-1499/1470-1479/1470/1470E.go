package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 1e18 + 7

type Branch struct {
	seq  []int
	pos  int
	cost int
	cnt  int64
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func buildDP(p []int, c int) [][]int64 {
	n := len(p) - 1 // p is 1-indexed
	dp := make([][]int64, n+2)
	for i := range dp {
		dp[i] = make([]int64, c+1)
	}
	for cost := 0; cost <= c; cost++ {
		dp[n+1][cost] = 1
	}
	for pos := n; pos >= 1; pos-- {
		for cost := 0; cost <= c; cost++ {
			val := dp[pos+1][cost] // no reversal
			for L := 2; L <= 5 && pos+L-1 <= n; L++ {
				if L-1 <= cost {
					val += dp[pos+L][cost-(L-1)]
					if val > INF {
						val = INF
					}
				}
			}
			dp[pos][cost] = val
		}
	}
	return dp
}

func findFirstGreater(dp [][]int64, start, n, cost int, j int64) int {
	lo := start
	hi := n + 1
	for lo < hi {
		mid := (lo + hi) / 2
		if dp[mid+1][cost] >= j {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func answer(p []int, dp [][]int64, n, c, idx int, j int64) int {
	if j > dp[1][c] {
		return -1
	}
	pos := 1
	cost := c
	curIdx := 1
	for pos <= n {
		// skip ranges without reversal
		if dp[pos+1][cost] >= j {
			// find farthest position where dp[next+1][cost] >= j
			next := findFirstGreater(dp, pos, n, cost, j)
			// positions pos..next-1 have no reversal
			if idx < curIdx+(next-pos) {
				return p[idx]
			}
			curIdx += next - pos
			pos = next
			if pos > n {
				break
			}
		} else {
			j -= dp[pos+1][cost]
			branches := make([]Branch, 0, 4)
			for L := 2; L <= 5 && pos+L-1 <= n; L++ {
				if L-1 <= cost {
					seq := make([]int, L)
					for k := 0; k < L; k++ {
						seq[k] = p[pos+L-1-k]
					}
					cnt := dp[pos+L][cost-(L-1)]
					if cnt > INF {
						cnt = INF
					}
					branches = append(branches, Branch{seq: seq, pos: pos + L, cost: cost - (L - 1), cnt: cnt})
				}
			}
			sort.Slice(branches, func(i, j int) bool {
				a := branches[i].seq
				b := branches[j].seq
				m := len(a)
				if len(b) < m {
					m = len(b)
				}
				for k := 0; k < m; k++ {
					if a[k] != b[k] {
						return a[k] < b[k]
					}
				}
				return len(a) < len(b)
			})
			for _, br := range branches {
				if j > br.cnt {
					j -= br.cnt
					continue
				}
				if idx-curIdx < len(br.seq) {
					return br.seq[idx-curIdx]
				}
				curIdx += len(br.seq)
				pos = br.pos
				cost = br.cost
				goto LOOP
			}
		}
	LOOP:
		if pos > n {
			break
		}
	}
	return p[idx]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, c, q int
		fmt.Fscan(reader, &n, &c, &q)
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		dp := buildDP(p, c)
		for ; q > 0; q-- {
			var i int
			var j int64
			fmt.Fscan(reader, &i, &j)
			ans := answer(p, dp, n, c, i, j)
			fmt.Fprintln(writer, ans)
		}
	}
}
