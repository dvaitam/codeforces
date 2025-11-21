package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

type operation struct {
	o int
	x int
}

func modPow(base int64, exp int64) int64 {
	res := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = res * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		ops := make([]operation, q)
		incA := make([]int, n)
		incB := make([]int, n)
		for i := 0; i < q; i++ {
			var o, x int
			fmt.Fscan(in, &o, &x)
			x--
			ops[i] = operation{o: o, x: x}
			if o == 1 {
				incA[x]++
			} else {
				incB[x]++
			}
		}

		values := make([]int64, 0, 2*n+q+5)
		for i := 0; i < n; i++ {
			for k := 0; k <= incA[i]; k++ {
				values = append(values, a[i]+int64(k))
			}
			for k := 0; k <= incB[i]; k++ {
				values = append(values, b[i]+int64(k))
			}
		}
		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
		coords := make([]int64, 0, len(values))
		for _, v := range values {
			if len(coords) == 0 || coords[len(coords)-1] != v {
				coords = append(coords, v)
			}
		}
		m := len(coords)
		index := make(map[int64]int, m)
		for i, v := range coords {
			index[v] = i
		}

		freqA := make([]int, m)
		freqB := make([]int, m)
		posA := make([]int, n)
		posB := make([]int, n)
		for i := 0; i < n; i++ {
			posA[i] = index[a[i]]
			freqA[posA[i]]++
			posB[i] = index[b[i]]
			freqB[posB[i]]++
		}

		cA := make([]int, m)
		cB := make([]int, m)
		running := 0
		for i := 0; i < m; i++ {
			running += freqA[i]
			cA[i] = running
		}
		running = 0
		for i := 0; i < m; i++ {
			running += freqB[i]
			cB[i] = running
		}

		kVals := make([]int, m)
		for i := 0; i < m; i++ {
			if cA[i] > cB[i] {
				kVals[i] = cA[i]
			} else {
				kVals[i] = cB[i]
			}
		}

		delta := make([]int, m)
		prev := 0
		for i := 0; i < m; i++ {
			delta[i] = kVals[i] - prev
			prev = kVals[i]
		}

		valMod := make([]int64, m)
		valInv := make([]int64, m)
		for i := 0; i < m; i++ {
			modVal := coords[i] % mod
			valMod[i] = modVal
			valInv[i] = modPow(modVal, mod-2)
		}

		ans := int64(1)
		for i := 0; i < m; i++ {
			if delta[i] == 0 {
				continue
			}
			ans = ans * modPow(valMod[i], int64(delta[i])) % mod
		}

		results := make([]int64, 0, q+1)
		results = append(results, ans)

		applyDelta := func(idx int, newDelta int) {
			if idx < 0 || idx >= m {
				return
			}
			diff := newDelta - delta[idx]
			if diff == 0 {
				return
			}
			if diff > 0 {
				ans = ans * modPow(valMod[idx], int64(diff)) % mod
			} else {
				ans = ans * modPow(valInv[idx], int64(-diff)) % mod
			}
			delta[idx] = newDelta
		}

		updateK := func(idx int) {
			newK := cA[idx]
			if cB[idx] > newK {
				newK = cB[idx]
			}
			if newK == kVals[idx] {
				return
			}
			kVals[idx] = newK
			prevK := 0
			if idx > 0 {
				prevK = kVals[idx-1]
			}
			applyDelta(idx, newK-prevK)
			if idx+1 < m {
				nextDelta := kVals[idx+1] - newK
				applyDelta(idx+1, nextDelta)
			}
		}

		for _, op := range ops {
			if op.o == 1 {
				idx := posA[op.x]
				cA[idx]--
				updateK(idx)
				posA[op.x]++
			} else {
				idx := posB[op.x]
				cB[idx]--
				updateK(idx)
				posB[op.x]++
			}
			results = append(results, ans)
		}

		for i, v := range results {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v%mod)
		}
		fmt.Fprintln(out)
	}
}
