package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

type BIT struct {
	n    int
	tree []int64
}

func newBIT(n int) *BIT {
	return &BIT{
		n:    n,
		tree: make([]int64, n+2),
	}
}

func (b *BIT) add(idx int, val int64) {
	for idx <= b.n {
		b.tree[idx] = (b.tree[idx] + val) % MOD
		if b.tree[idx] < 0 {
			b.tree[idx] += MOD
		}
		idx += idx & -idx
	}
}

func (b *BIT) sum(idx int) int64 {
	if idx <= 0 {
		return 0
	}
	if idx > b.n {
		idx = b.n
	}
	var res int64
	for idx > 0 {
		res += b.tree[idx]
		if res >= MOD {
			res -= MOD
		}
		idx -= idx & -idx
	}
	return res
}

func (b *BIT) rangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return (b.sum(r) - b.sum(l-1) + MOD) % MOD
}

var pow2 []int64

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const maxN = 200000 + 5
	pow2 = make([]int64, maxN)
	pow2[0] = 1
	for i := 1; i < maxN; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		bytes := []byte(s)

		bitPref0 := newBIT(n)
		bitPref1 := newBIT(n)
		bitSuf0 := newBIT(n)
		bitSuf1 := newBIT(n)

		var trans int64
		var pref0, pref1 int64
		for i := 1; i <= n; i++ {
			valPref := pow2[i-1]
			valSuf := pow2[n-i]
			if bytes[i-1] == '0' {
				bitPref0.add(i, valPref)
				bitSuf0.add(i, valSuf)
				trans = (trans + valSuf*pref1) % MOD
				pref0 = (pref0 + valPref) % MOD
			} else {
				bitPref1.add(i, valPref)
				bitSuf1.add(i, valSuf)
				trans = (trans + valSuf*pref0) % MOD
				pref1 = (pref1 + valPref) % MOD
			}
		}

		var q int
		fmt.Fscan(in, &q)
		arr := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &arr[i])
		}

		base := (pow2[n] - 1 + MOD) % MOD

		for _, pos := range arr {
			oldChar := bytes[pos-1]
			var prefOpp, sufOpp, valPref, valSuf int64
			valPref = pow2[pos-1]
			valSuf = pow2[n-pos]

			if oldChar == '0' {
				prefOpp = bitPref1.sum(pos - 1)
				sufOpp = bitSuf1.rangeSum(pos+1, n)
				contrib := (valSuf*prefOpp + valPref*sufOpp) % MOD
				trans = (trans - contrib) % MOD
				if trans < 0 {
					trans += MOD
				}
				bitPref0.add(pos, -valPref)
				bitSuf0.add(pos, -valSuf)
			} else {
				prefOpp = bitPref0.sum(pos - 1)
				sufOpp = bitSuf0.rangeSum(pos+1, n)
				contrib := (valSuf*prefOpp + valPref*sufOpp) % MOD
				trans = (trans - contrib) % MOD
				if trans < 0 {
					trans += MOD
				}
				bitPref1.add(pos, -valPref)
				bitSuf1.add(pos, -valSuf)
			}

			var newChar byte
			if oldChar == '0' {
				newChar = '1'
			} else {
				newChar = '0'
			}
			bytes[pos-1] = newChar

			if newChar == '0' {
				prefOpp = bitPref1.sum(pos - 1)
				sufOpp = bitSuf1.rangeSum(pos+1, n)
				contrib := (valSuf*prefOpp + valPref*sufOpp) % MOD
				trans = (trans + contrib) % MOD
				bitPref0.add(pos, valPref)
				bitSuf0.add(pos, valSuf)
			} else {
				prefOpp = bitPref0.sum(pos - 1)
				sufOpp = bitSuf0.rangeSum(pos+1, n)
				contrib := (valSuf*prefOpp + valPref*sufOpp) % MOD
				trans = (trans + contrib) % MOD
				bitPref1.add(pos, valPref)
				bitSuf1.add(pos, valSuf)
			}

			ans := (base + trans) % MOD
			fmt.Fprintf(out, "%d ", ans)
		}
		fmt.Fprintln(out)
	}
}
