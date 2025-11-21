package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type hashPair struct {
	x uint64
	y uint64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	rng := rand.New(rand.NewSource(123456789))

	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		hashes := make(map[int]hashPair)
		prefX := make([]uint64, n+1)
		prefY := make([]uint64, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if _, ok := hashes[a[i]]; !ok {
				hashes[a[i]] = hashPair{rng.Uint64(), rng.Uint64()}
			}
			h := hashes[a[i]]
			prefX[i+1] = prefX[i] ^ h.x
			prefY[i+1] = prefY[i] ^ h.y
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if (r-l+1)%2 == 1 {
				fmt.Fprintln(out, "NO")
				continue
			}
			lx := prefX[r] ^ prefX[l-1]
			ly := prefY[r] ^ prefY[l-1]
			if lx == 0 && ly == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
