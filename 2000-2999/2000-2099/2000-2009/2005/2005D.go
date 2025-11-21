package main

import (
	"bufio"
	"fmt"
	"os"
)

type entry struct {
	gA    int64
	gB    int64
	start int
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}

		prefA := make([]int64, n+1)
		prefB := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefA[i] = gcd(prefA[i-1], a[i])
			prefB[i] = gcd(prefB[i-1], b[i])
		}
		sufA := make([]int64, n+3)
		sufB := make([]int64, n+3)
		for i := n; i >= 1; i-- {
			sufA[i] = gcd(sufA[i+1], a[i])
			sufB[i] = gcd(sufB[i+1], b[i])
		}

		prefBlockA := make([]int, n+1)
		prefBlockB := make([]int, n+1)
		prefBlockA[0] = 1
		prefBlockB[0] = 1
		for i := 1; i <= n; i++ {
			if prefA[i] == prefA[i-1] {
				prefBlockA[i] = prefBlockA[i-1]
			} else {
				prefBlockA[i] = i + 1
			}
			if prefB[i] == prefB[i-1] {
				prefBlockB[i] = prefBlockB[i-1]
			} else {
				prefBlockB[i] = i + 1
			}
		}

		best := int64(-1)
		count := int64(0)
		vec := make([]entry, 0)

		for r := 1; r <= n; r++ {
			newVec := make([]entry, 0, len(vec)+1)
			newVec = append(newVec, entry{a[r], b[r], r})
			for _, e := range vec {
				gA := gcd(e.gA, a[r])
				gB := gcd(e.gB, b[r])
				last := &newVec[len(newVec)-1]
				if gA == last.gA && gB == last.gB {
					last.start = e.start
				} else {
					newVec = append(newVec, entry{gA, gB, e.start})
				}
			}
			vec = newVec

			sA := sufA[r+1]
			sB := sufB[r+1]

			prevStart := r + 1
			for _, e := range vec {
				L := e.start
				R := prevStart - 1
				if L > R {
					prevStart = L
					continue
				}
				gAinside := e.gA
				gBinside := e.gB
				lCur := R
				for lCur >= L {
					pIdx := lCur - 1
					limitA := prefBlockA[pIdx]
					if limitA < L {
						limitA = L
					}
					limitB := prefBlockB[pIdx]
					if limitB < L {
						limitB = L
					}
					blockL := limitA
					if limitB > blockL {
						blockL = limitB
					}
					pAVal := prefA[pIdx]
					pBVal := prefB[pIdx]
					baseA := gcd(pAVal, sA)
					baseB := gcd(pBVal, sB)
					gAres := gcd(baseA, gBinside)
					gBres := gcd(baseB, gAinside)
					total := gAres + gBres
					length := int64(lCur - blockL + 1)
					if total > best {
						best = total
						count = length
					} else if total == best {
						count += length
					}
					lCur = blockL - 1
				}
				prevStart = L
			}
		}

		fmt.Fprintln(out, best, count)
	}
}
