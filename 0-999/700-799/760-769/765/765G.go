package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	m := len(s)

	var n int
	fmt.Fscan(in, &n)
	type PA struct {
		p int
		a int
	}
	arr := make([]PA, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i].p, &arr[i].a)
	}

	// collect zero and one positions
	zeroPos := make([]int, 0)
	idxMap := make([]int, m)
	for i := 0; i < m; i++ {
		if s[i] == '0' {
			idxMap[i] = len(zeroPos)
			zeroPos = append(zeroPos, i)
		} else {
			idxMap[i] = -1
		}
	}
	z := len(zeroPos)
	fullMask := uint64(1<<z) - 1

	dp := map[uint64]*big.Int{0: big.NewInt(1)}

	for _, pa := range arr {
		p := pa.p
		a := pa.a
		pow := new(big.Int).Exp(big.NewInt(int64(p)), big.NewInt(int64(a-1)), nil)

		patterns := make(map[uint64]int64)
		if p > m {
			patterns[0] = int64(p - m)
			for i := 0; i < m; i++ {
				if s[i] == '0' {
					mask := uint64(1) << uint(idxMap[i])
					patterns[mask]++
				}
			}
		} else {
			used := make([]bool, p)
			masks := make([]uint64, p)
			valid := make([]bool, p)
			for i := range valid {
				valid[i] = true
			}
			for i := 0; i < m; i++ {
				r := (p - (i % p)) % p
				used[r] = true
				if s[i] == '1' {
					valid[r] = false
				} else {
					masks[r] |= 1 << uint(idxMap[i])
				}
			}
			for r := 0; r < p; r++ {
				if !used[r] {
					patterns[0]++
				} else if valid[r] {
					patterns[masks[r]]++
				}
			}
		}

		newdp := make(map[uint64]*big.Int)
		for mask, val := range dp {
			for pat, cnt := range patterns {
				if cnt == 0 {
					continue
				}
				nm := mask | pat
				add := new(big.Int).Mul(val, pow)
				add.Mul(add, big.NewInt(cnt))
				if ex, ok := newdp[nm]; ok {
					ex.Add(ex, add)
				} else {
					newdp[nm] = add
				}
			}
		}
		dp = newdp
	}

	ans := dp[fullMask]
	if ans == nil {
		ans = new(big.Int)
	}
	fmt.Fprintln(out, ans.String())
}
