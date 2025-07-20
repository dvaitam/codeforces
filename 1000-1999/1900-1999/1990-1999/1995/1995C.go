package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for tt := 0;tt < t;tt++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		total := int64(0)
		prev_base := a[0]
		prev_k := 0
		ok := true
		for i := 1; i < n && ok; i++ {
			curr_a := a[i]
			var curr_k int
			if curr_a == 1 {
				if prev_base != 1 {
					ok = false
					break
				}
				curr_k = 0
			} else {
				min_k := math.MaxInt
				found := false
				start := max(0, prev_k-6)
				end := prev_k + 6
				for ck := start; ck <= end; ck++ {
					if isGE(curr_a, ck, prev_base, prev_k) {
						if ck < min_k {
							min_k = ck
						}
						found = true
					}
				}
				if !found {
					ok = false
					break
				}
				curr_k = min_k
			}
			total += int64(curr_k)
			prev_base = curr_a
			prev_k = curr_k
		}
		if !ok {
			fmt.Println(-1)
		} else {
			fmt.Println(total)
		}
	}
}

func isGE(a int64, k int, pb int64, pk int) bool {
	if k == pk {
		return a >= pb
	}
	if k > pk {
		d := k - pk
		exp := int64(1 << uint(d))
		powa := new(big.Int).Exp(big.NewInt(a), big.NewInt(exp), nil)
		return powa.Cmp(big.NewInt(pb)) >= 0
	} else {
		d := pk - k
		exp := int64(1 << uint(d))
		powpb := new(big.Int).Exp(big.NewInt(pb), big.NewInt(exp), nil)
		return big.NewInt(a).Cmp(powpb) >= 0
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
