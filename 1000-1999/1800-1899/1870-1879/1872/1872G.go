package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n+1)
		prefix := make([]int64, n+1)
		idx := make([]int, 0)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			prefix[i] = prefix[i-1] + a[i]
			if a[i] > 1 {
				idx = append(idx, i)
			}
		}
		if len(idx) == 0 {
			fmt.Fprintln(writer, 1, 1)
			continue
		}
		if len(idx) > 60 {
			fmt.Fprintln(writer, idx[0], idx[len(idx)-1])
			continue
		}
		bestL, bestR := idx[0], idx[0]
		bestDiff := big.NewInt(0)
		for i := 0; i < len(idx); i++ {
			prod := big.NewInt(1)
			l := idx[i]
			for j := i; j < len(idx); j++ {
				prod.Mul(prod, big.NewInt(a[idx[j]]))
				r := idx[j]
				sumSub := prefix[r] - prefix[l-1]
				diff := new(big.Int).Sub(prod, big.NewInt(sumSub))
				if diff.Cmp(bestDiff) > 0 {
					bestDiff = diff
					bestL = l
					bestR = r
				}
			}
		}
		fmt.Fprintln(writer, bestL, bestR)
	}
}
