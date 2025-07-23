package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var l, wmax float64
	if _, err := fmt.Fscan(in, &n, &l, &wmax); err != nil {
		return
	}
	xs := make([]float64, n)
	vs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &vs[i])
	}
	// sort by position just in case input is not ordered
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return xs[idx[i]] < xs[idx[j]] })

	ans := 0
	for ii := 0; ii < n; ii++ {
		i := idx[ii]
		for jj := ii + 1; jj < n; jj++ {
			j := idx[jj]
			if vs[i] == vs[j] {
				continue
			}
			vi := float64(vs[i])
			vj := float64(vs[j])
			denom := xs[i] - xs[j]
			if denom == 0 {
				continue
			}
			w := ((xs[j]+l/2)*vi - (xs[i]+l/2)*vj) / denom
			if math.Abs(w) > wmax {
				continue
			}
			di := vi + w
			dj := vj + w
			if di == 0 || dj == 0 {
				continue
			}
			t := -(xs[i] + l/2) / di
			if t <= 0 {
				continue
			}
			// if centers coincide, clouds surely cover the moon at that moment
			ans++
		}
	}
	fmt.Println(ans)
}
