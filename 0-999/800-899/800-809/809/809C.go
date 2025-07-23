package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

var A, B, limitK int64
var memo [31][2][2][2]struct {
	cnt int64
	sum int64
}
var used [31][2][2][2]bool

func dfs(pos int, la, lb, lk bool) (int64, int64) {
	if pos < 0 {
		return 1, 0
	}
	iLa, iLb, iLk := 0, 0, 0
	if la {
		iLa = 1
	}
	if lb {
		iLb = 1
	}
	if lk {
		iLk = 1
	}
	if used[pos][iLa][iLb][iLk] {
		res := memo[pos][iLa][iLb][iLk]
		return res.cnt, res.sum
	}
	bitA := (A >> pos) & 1
	bitB := (B >> pos) & 1
	bitK := (limitK >> pos) & 1
	var resCnt, resSum int64
	for da := int64(0); da <= 1; da++ {
		if la && da > bitA {
			continue
		}
		nLa := la && da == bitA
		for db := int64(0); db <= 1; db++ {
			if lb && db > bitB {
				continue
			}
			nLb := lb && db == bitB
			xr := da ^ db
			if lk && xr > bitK {
				continue
			}
			nLk := lk && xr == bitK
			c, s := dfs(pos-1, nLa, nLb, nLk)
			resCnt += c
			resSum = (resSum + s + (c%mod)*((xr<<pos)%mod)) % mod
		}
	}
	used[pos][iLa][iLb][iLk] = true
	memo[pos][iLa][iLb][iLk] = struct{ cnt, sum int64 }{resCnt, resSum}
	return resCnt, resSum
}

func calc(a, b, k int64) int64 {
	if a < 0 || b < 0 || k <= 0 {
		return 0
	}
	A, B = a, b
	limitK = k - 1
	for i := range used {
		for j := range used[i] {
			for l := range used[i][j] {
				used[i][j][l][0] = false
				used[i][j][l][1] = false
			}
		}
	}
	cnt, sum := dfs(30, true, true, true)
	return (sum + cnt) % mod
}

func query(x1, y1, x2, y2, k int64) int64 {
	res := calc(x2-1, y2-1, k)
	res -= calc(x1-2, y2-1, k)
	res -= calc(x2-1, y1-2, k)
	res += calc(x1-2, y1-2, k)
	res %= mod
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var x1, y1, x2, y2, k int64
		fmt.Fscan(in, &x1, &y1, &x2, &y2, &k)
		ans := query(x1, y1, x2, y2, k)
		fmt.Fprintln(out, ans)
	}
}
