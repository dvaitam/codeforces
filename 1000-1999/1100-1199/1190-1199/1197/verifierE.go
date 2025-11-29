package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const solution1197ESource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD = 1000000007
const INF = int64(4e18)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	dolls := make([]struct{ out, inV int }, n)
	outs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &dolls[i].out, &dolls[i].inV)
		outs[i] = dolls[i].out
	}
	Omin := outs[0]
	Imax := dolls[0].inV
	for i := 1; i < n; i++ {
		if dolls[i].out < Omin {
			Omin = dolls[i].out
		}
		if dolls[i].inV > Imax {
			Imax = dolls[i].inV
		}
	}
	sort.Ints(outs)
	uouts := outs[:1]
	for i := 1; i < n; i++ {
		if outs[i] != outs[i-1] {
			uouts = append(uouts, outs[i])
		}
	}
	m := len(uouts)
	fcoord := make([]int64, m)
	ccoord := make([]int, m)
	for i := 0; i < m; i++ {
		fcoord[i] = INF
		ccoord[i] = 0
	}
	sort.Slice(dolls, func(i, j int) bool {
		if dolls[i].inV != dolls[j].inV {
			return dolls[i].inV < dolls[j].inV
		}
		return dolls[i].out < dolls[j].out
	})
	dp := make([]int64, n)
	cnt := make([]int, n)
	for i, d := range dolls {
		dp[i] = INF
		cnt[i] = 0
		if d.out == Omin {
			dp[i] = int64(d.inV)
			cnt[i] = 1
		}
		r := sort.Search(len(uouts), func(j int) bool { return uouts[j] > d.inV })
		if r > 0 {
			if fcoord[r-1] < INF {
				val := fcoord[r-1] + int64(d.inV)
				if val < dp[i] {
					dp[i] = val
					cnt[i] = ccoord[r-1]
				} else if val == dp[i] {
					cnt[i] = (cnt[i] + ccoord[r-1]) % MOD
				}
			}
		}
		if dp[i] < INF {
			pos := sort.Search(len(uouts), func(j int) bool { return uouts[j] >= d.out })
			fnew := dp[i] - int64(d.out)
			if fnew < fcoord[pos] {
				fcoord[pos] = fnew
				ccoord[pos] = cnt[i]
			} else if fnew == fcoord[pos] {
				ccoord[pos] = (ccoord[pos] + cnt[i]) % MOD
			}
		}
	}
	best := INF
	ways := 0
	for i, d := range dolls {
		if d.inV == Imax && dp[i] < INF {
			if dp[i] < best {
				best = dp[i]
				ways = cnt[i]
			} else if dp[i] == best {
				ways = (ways + cnt[i]) % MOD
			}
		}
	}
	fmt.Fprint(out, ways)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1197ESource

type pair struct {
	out int
	inV int
}

type testCase struct {
	n     int
	dolls []pair
}

var testcases = []testCase{
	{n: 2, dolls: []pair{{out: 11, inV: 2}, {out: 14, inV: 8}}},
	{n: 2, dolls: []pair{{out: 4, inV: 1}, {out: 2, inV: 1}}},
	{n: 5, dolls: []pair{{out: 11, inV: 1}, {out: 9, inV: 6}, {out: 10, inV: 3}, {out: 5, inV: 3}, {out: 8, inV: 1}}},
	{n: 6, dolls: []pair{{out: 10, inV: 5}, {out: 8, inV: 2}, {out: 11, inV: 5}, {out: 13, inV: 2}, {out: 12, inV: 11}, {out: 14, inV: 9}}},
	{n: 2, dolls: []pair{{out: 7, inV: 2}, {out: 17, inV: 9}}},
	{n: 1, dolls: []pair{{out: 19, inV: 10}}},
	{n: 1, dolls: []pair{{out: 11, inV: 10}}},
	{n: 6, dolls: []pair{{out: 11, inV: 9}, {out: 8, inV: 4}, {out: 15, inV: 10}, {out: 11, inV: 7}, {out: 16, inV: 3}, {out: 9, inV: 5}}},
	{n: 3, dolls: []pair{{out: 3, inV: 1}, {out: 3, inV: 2}, {out: 10, inV: 9}}},
	{n: 5, dolls: []pair{{out: 17, inV: 11}, {out: 6, inV: 2}, {out: 4, inV: 2}, {out: 8, inV: 6}, {out: 16, inV: 5}}},
	{n: 2, dolls: []pair{{out: 13, inV: 7}, {out: 20, inV: 11}}},
	{n: 6, dolls: []pair{{out: 19, inV: 7}, {out: 12, inV: 2}, {out: 3, inV: 1}, {out: 10, inV: 4}, {out: 5, inV: 3}, {out: 7, inV: 3}}},
	{n: 4, dolls: []pair{{out: 2, inV: 1}, {out: 13, inV: 12}, {out: 4, inV: 2}, {out: 12, inV: 1}}},
	{n: 3, dolls: []pair{{out: 11, inV: 6}, {out: 6, inV: 4}, {out: 4, inV: 2}}},
	{n: 5, dolls: []pair{{out: 8, inV: 4}, {out: 11, inV: 3}, {out: 10, inV: 7}, {out: 7, inV: 3}, {out: 20, inV: 1}}},
	{n: 3, dolls: []pair{{out: 3, inV: 2}, {out: 7, inV: 3}, {out: 13, inV: 5}}},
	{n: 5, dolls: []pair{{out: 5, inV: 4}, {out: 8, inV: 4}, {out: 8, inV: 1}, {out: 3, inV: 1}, {out: 3, inV: 1}}},
	{n: 5, dolls: []pair{{out: 6, inV: 5}, {out: 3, inV: 2}, {out: 20, inV: 8}, {out: 12, inV: 1}, {out: 5, inV: 3}}},
	{n: 4, dolls: []pair{{out: 8, inV: 4}, {out: 8, inV: 2}, {out: 16, inV: 7}, {out: 17, inV: 2}}},
	{n: 2, dolls: []pair{{out: 15, inV: 8}, {out: 9, inV: 7}}},
	{n: 2, dolls: []pair{{out: 17, inV: 7}, {out: 3, inV: 1}}},
	{n: 3, dolls: []pair{{out: 10, inV: 4}, {out: 18, inV: 7}, {out: 9, inV: 7}}},
	{n: 3, dolls: []pair{{out: 6, inV: 3}, {out: 3, inV: 2}, {out: 20, inV: 4}}},
	{n: 5, dolls: []pair{{out: 14, inV: 11}, {out: 3, inV: 2}, {out: 14, inV: 2}, {out: 15, inV: 4}, {out: 20, inV: 6}}},
	{n: 3, dolls: []pair{{out: 11, inV: 8}, {out: 12, inV: 7}, {out: 18, inV: 7}}},
	{n: 6, dolls: []pair{{out: 10, inV: 6}, {out: 14, inV: 8}, {out: 4, inV: 2}, {out: 8, inV: 1}, {out: 14, inV: 10}, {out: 6, inV: 3}}},
	{n: 6, dolls: []pair{{out: 3, inV: 1}, {out: 16, inV: 10}, {out: 17, inV: 13}, {out: 14, inV: 4}, {out: 2, inV: 1}, {out: 7, inV: 1}}},
	{n: 5, dolls: []pair{{out: 10, inV: 2}, {out: 14, inV: 13}, {out: 14, inV: 4}, {out: 19, inV: 2}, {out: 8, inV: 2}}},
	{n: 6, dolls: []pair{{out: 12, inV: 9}, {out: 17, inV: 15}, {out: 2, inV: 1}, {out: 3, inV: 1}, {out: 17, inV: 9}, {out: 6, inV: 1}}},
	{n: 3, dolls: []pair{{out: 4, inV: 3}, {out: 2, inV: 1}, {out: 13, inV: 2}}},
	{n: 1, dolls: []pair{{out: 19, inV: 15}}},
	{n: 4, dolls: []pair{{out: 8, inV: 7}, {out: 11, inV: 7}, {out: 9, inV: 8}, {out: 14, inV: 2}}},
	{n: 1, dolls: []pair{{out: 5, inV: 3}}},
	{n: 5, dolls: []pair{{out: 15, inV: 7}, {out: 16, inV: 2}, {out: 8, inV: 6}, {out: 11, inV: 8}, {out: 15, inV: 2}}},
	{n: 5, dolls: []pair{{out: 7, inV: 3}, {out: 7, inV: 2}, {out: 6, inV: 3}, {out: 17, inV: 11}, {out: 10, inV: 9}}},
	{n: 1, dolls: []pair{{out: 7, inV: 1}}},
	{n: 6, dolls: []pair{{out: 11, inV: 2}, {out: 19, inV: 4}, {out: 17, inV: 16}, {out: 18, inV: 3}, {out: 18, inV: 8}, {out: 15, inV: 14}}},
	{n: 3, dolls: []pair{{out: 13, inV: 4}, {out: 7, inV: 6}, {out: 2, inV: 1}}},
	{n: 5, dolls: []pair{{out: 12, inV: 9}, {out: 16, inV: 13}, {out: 20, inV: 10}, {out: 18, inV: 15}, {out: 16, inV: 7}}},
	{n: 2, dolls: []pair{{out: 10, inV: 6}, {out: 12, inV: 3}}},
	{n: 4, dolls: []pair{{out: 4, inV: 3}, {out: 6, inV: 5}, {out: 7, inV: 3}, {out: 13, inV: 4}}},
	{n: 5, dolls: []pair{{out: 13, inV: 11}, {out: 4, inV: 1}, {out: 14, inV: 11}, {out: 7, inV: 3}, {out: 13, inV: 6}}},
	{n: 2, dolls: []pair{{out: 11, inV: 1}, {out: 2, inV: 1}}},
	{n: 3, dolls: []pair{{out: 5, inV: 2}, {out: 7, inV: 5}, {out: 17, inV: 3}}},
	{n: 1, dolls: []pair{{out: 7, inV: 6}}},
	{n: 4, dolls: []pair{{out: 9, inV: 5}, {out: 14, inV: 10}, {out: 9, inV: 8}, {out: 9, inV: 5}}},
	{n: 3, dolls: []pair{{out: 9, inV: 6}, {out: 19, inV: 17}, {out: 16, inV: 15}}},
	{n: 4, dolls: []pair{{out: 18, inV: 13}, {out: 12, inV: 5}, {out: 16, inV: 7}, {out: 20, inV: 1}}},
	{n: 3, dolls: []pair{{out: 7, inV: 5}, {out: 16, inV: 12}, {out: 19, inV: 12}}},
	{n: 4, dolls: []pair{{out: 14, inV: 10}, {out: 2, inV: 1}, {out: 18, inV: 3}, {out: 16, inV: 11}}},
	{n: 3, dolls: []pair{{out: 11, inV: 2}, {out: 10, inV: 8}, {out: 9, inV: 8}}},
	{n: 5, dolls: []pair{{out: 4, inV: 1}, {out: 9, inV: 2}, {out: 11, inV: 3}, {out: 3, inV: 1}, {out: 14, inV: 13}}},
	{n: 5, dolls: []pair{{out: 19, inV: 18}, {out: 10, inV: 1}, {out: 18, inV: 8}, {out: 7, inV: 1}, {out: 8, inV: 1}}},
	{n: 3, dolls: []pair{{out: 4, inV: 1}, {out: 10, inV: 1}, {out: 11, inV: 10}}},
	{n: 6, dolls: []pair{{out: 7, inV: 6}, {out: 6, inV: 4}, {out: 6, inV: 1}, {out: 19, inV: 12}, {out: 2, inV: 1}, {out: 18, inV: 14}}},
	{n: 2, dolls: []pair{{out: 8, inV: 7}, {out: 11, inV: 8}}},
	{n: 5, dolls: []pair{{out: 4, inV: 2}, {out: 7, inV: 2}, {out: 10, inV: 9}, {out: 14, inV: 9}, {out: 11, inV: 7}}},
	{n: 3, dolls: []pair{{out: 7, inV: 4}, {out: 3, inV: 2}, {out: 2, inV: 1}}},
	{n: 6, dolls: []pair{{out: 2, inV: 1}, {out: 6, inV: 1}, {out: 7, inV: 1}, {out: 2, inV: 1}, {out: 9, inV: 1}, {out: 17, inV: 6}}},
	{n: 4, dolls: []pair{{out: 14, inV: 6}, {out: 7, inV: 5}, {out: 20, inV: 7}, {out: 5, inV: 4}}},
	{n: 5, dolls: []pair{{out: 13, inV: 5}, {out: 15, inV: 10}, {out: 3, inV: 1}, {out: 10, inV: 5}, {out: 17, inV: 11}}},
	{n: 1, dolls: []pair{{out: 9, inV: 6}}},
	{n: 1, dolls: []pair{{out: 3, inV: 2}}},
	{n: 5, dolls: []pair{{out: 17, inV: 12}, {out: 4, inV: 1}, {out: 3, inV: 2}, {out: 18, inV: 8}, {out: 3, inV: 1}}},
	{n: 6, dolls: []pair{{out: 4, inV: 2}, {out: 6, inV: 3}, {out: 5, inV: 1}, {out: 4, inV: 1}, {out: 10, inV: 1}, {out: 3, inV: 1}}},
	{n: 4, dolls: []pair{{out: 6, inV: 2}, {out: 13, inV: 4}, {out: 13, inV: 5}, {out: 14, inV: 11}}},
	{n: 5, dolls: []pair{{out: 14, inV: 1}, {out: 7, inV: 4}, {out: 17, inV: 3}, {out: 18, inV: 17}, {out: 11, inV: 9}}},
	{n: 6, dolls: []pair{{out: 11, inV: 4}, {out: 15, inV: 6}, {out: 7, inV: 2}, {out: 2, inV: 1}, {out: 15, inV: 5}, {out: 12, inV: 7}}},
	{n: 6, dolls: []pair{{out: 15, inV: 6}, {out: 8, inV: 3}, {out: 10, inV: 9}, {out: 4, inV: 1}, {out: 14, inV: 5}, {out: 11, inV: 3}}},
	{n: 5, dolls: []pair{{out: 12, inV: 2}, {out: 7, inV: 5}, {out: 4, inV: 3}, {out: 8, inV: 3}, {out: 8, inV: 6}}},
	{n: 4, dolls: []pair{{out: 5, inV: 3}, {out: 2, inV: 1}, {out: 20, inV: 7}, {out: 7, inV: 1}}},
	{n: 2, dolls: []pair{{out: 18, inV: 4}, {out: 5, inV: 3}}},
	{n: 3, dolls: []pair{{out: 20, inV: 19}, {out: 6, inV: 1}, {out: 7, inV: 6}}},
	{n: 6, dolls: []pair{{out: 2, inV: 1}, {out: 8, inV: 3}, {out: 3, inV: 2}, {out: 12, inV: 10}, {out: 13, inV: 3}, {out: 7, inV: 2}}},
	{n: 1, dolls: []pair{{out: 3, inV: 1}}},
	{n: 6, dolls: []pair{{out: 3, inV: 2}, {out: 17, inV: 4}, {out: 8, inV: 7}, {out: 6, inV: 4}, {out: 17, inV: 16}, {out: 19, inV: 11}}},
	{n: 3, dolls: []pair{{out: 15, inV: 12}, {out: 14, inV: 1}, {out: 3, inV: 1}}},
	{n: 2, dolls: []pair{{out: 13, inV: 4}, {out: 11, inV: 4}}},
	{n: 3, dolls: []pair{{out: 20, inV: 7}, {out: 13, inV: 7}, {out: 13, inV: 9}}},
	{n: 6, dolls: []pair{{out: 5, inV: 2}, {out: 8, inV: 3}, {out: 7, inV: 5}, {out: 8, inV: 1}, {out: 2, inV: 1}, {out: 8, inV: 3}}},
	{n: 6, dolls: []pair{{out: 12, inV: 9}, {out: 17, inV: 2}, {out: 11, inV: 7}, {out: 7, inV: 2}, {out: 14, inV: 10}, {out: 16, inV: 6}}},
	{n: 3, dolls: []pair{{out: 2, inV: 1}, {out: 6, inV: 4}, {out: 2, inV: 1}}},
	{n: 2, dolls: []pair{{out: 5, inV: 3}, {out: 11, inV: 10}}},
	{n: 5, dolls: []pair{{out: 3, inV: 1}, {out: 17, inV: 6}, {out: 18, inV: 2}, {out: 7, inV: 5}, {out: 9, inV: 6}}},
	{n: 5, dolls: []pair{{out: 2, inV: 1}, {out: 14, inV: 7}, {out: 9, inV: 3}, {out: 3, inV: 1}, {out: 7, inV: 5}}},
	{n: 3, dolls: []pair{{out: 9, inV: 7}, {out: 12, inV: 9}, {out: 13, inV: 3}}},
	{n: 6, dolls: []pair{{out: 13, inV: 11}, {out: 19, inV: 12}, {out: 10, inV: 6}, {out: 19, inV: 2}, {out: 19, inV: 5}, {out: 13, inV: 10}}},
	{n: 2, dolls: []pair{{out: 17, inV: 16}, {out: 5, inV: 1}}},
	{n: 6, dolls: []pair{{out: 8, inV: 2}, {out: 6, inV: 5}, {out: 6, inV: 5}, {out: 3, inV: 2}, {out: 4, inV: 3}, {out: 3, inV: 1}}},
	{n: 5, dolls: []pair{{out: 7, inV: 3}, {out: 3, inV: 2}, {out: 17, inV: 9}, {out: 8, inV: 6}, {out: 17, inV: 10}}},
	{n: 5, dolls: []pair{{out: 3, inV: 1}, {out: 11, inV: 7}, {out: 3, inV: 1}, {out: 13, inV: 9}, {out: 13, inV: 2}}},
	{n: 2, dolls: []pair{{out: 15, inV: 11}, {out: 4, inV: 2}}},
	{n: 1, dolls: []pair{{out: 11, inV: 4}}},
	{n: 6, dolls: []pair{{out: 8, inV: 7}, {out: 7, inV: 2}, {out: 4, inV: 1}, {out: 2, inV: 1}, {out: 11, inV: 9}, {out: 18, inV: 17}}},
	{n: 4, dolls: []pair{{out: 4, inV: 1}, {out: 18, inV: 11}, {out: 7, inV: 6}, {out: 8, inV: 5}}},
	{n: 6, dolls: []pair{{out: 19, inV: 2}, {out: 3, inV: 1}, {out: 6, inV: 4}, {out: 6, inV: 1}, {out: 18, inV: 11}, {out: 18, inV: 11}}},
	{n: 6, dolls: []pair{{out: 10, inV: 3}, {out: 20, inV: 19}, {out: 5, inV: 1}, {out: 2, inV: 1}, {out: 3, inV: 2}, {out: 19, inV: 14}}},
	{n: 3, dolls: []pair{{out: 20, inV: 12}, {out: 6, inV: 5}, {out: 13, inV: 9}}},
	{n: 4, dolls: []pair{{out: 9, inV: 8}, {out: 7, inV: 4}, {out: 20, inV: 18}, {out: 2, inV: 1}}},
	{n: 3, dolls: []pair{{out: 20, inV: 19}, {out: 16, inV: 6}, {out: 7, inV: 3}}},
}

const (
	mod = 1000000007
	inf = int64(4e18)
)

func solveCase(tc testCase) string {
	n := tc.n
	dolls := make([]pair, n)
	outs := make([]int, n)
	for i, d := range tc.dolls {
		dolls[i] = d
		outs[i] = d.out
	}
	Omin := outs[0]
	Imax := dolls[0].inV
	for i := 1; i < n; i++ {
		if dolls[i].out < Omin {
			Omin = dolls[i].out
		}
		if dolls[i].inV > Imax {
			Imax = dolls[i].inV
		}
	}
	sort.Ints(outs)
	uouts := outs[:1]
	for i := 1; i < n; i++ {
		if outs[i] != outs[i-1] {
			uouts = append(uouts, outs[i])
		}
	}
	m := len(uouts)
	fcoord := make([]int64, m)
	ccoord := make([]int, m)
	for i := 0; i < m; i++ {
		fcoord[i] = inf
		ccoord[i] = 0
	}
	sort.Slice(dolls, func(i, j int) bool {
		if dolls[i].inV != dolls[j].inV {
			return dolls[i].inV < dolls[j].inV
		}
		return dolls[i].out < dolls[j].out
	})
	dp := make([]int64, n)
	cnt := make([]int, n)
	for i, d := range dolls {
		dp[i] = inf
		if d.out == Omin {
			dp[i] = int64(d.inV)
			cnt[i] = 1
		}
		r := sort.Search(len(uouts), func(j int) bool { return uouts[j] > d.inV })
		if r > 0 && fcoord[r-1] < inf {
			val := fcoord[r-1] + int64(d.inV)
			if val < dp[i] {
				dp[i] = val
				cnt[i] = ccoord[r-1]
			} else if val == dp[i] {
				cnt[i] = (cnt[i] + ccoord[r-1]) % mod
			}
		}
		if dp[i] < inf {
			pos := sort.Search(len(uouts), func(j int) bool { return uouts[j] >= d.out })
			fnew := dp[i] - int64(d.out)
			if fnew < fcoord[pos] {
				fcoord[pos] = fnew
				ccoord[pos] = cnt[i]
			} else if fnew == fcoord[pos] {
				ccoord[pos] = (ccoord[pos] + cnt[i]) % mod
			}
		}
	}
	best := inf
	ways := 0
	for i, d := range dolls {
		if d.inV == Imax && dp[i] < inf {
			if dp[i] < best {
				best = dp[i]
				ways = cnt[i]
			} else if dp[i] == best {
				ways = (ways + cnt[i]) % mod
			}
		}
	}
	return fmt.Sprintf("%d", ways)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, d := range tc.dolls {
			fmt.Fprintf(&sb, "%d %d\n", d.out, d.inV)
		}
		input := sb.String()
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
