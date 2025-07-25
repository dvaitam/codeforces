package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	t int64
	a int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		x := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}
		pairs := make([]pair, n)
		for i := 0; i < n; i++ {
			t := x[i]
			if t < 0 {
				t = -t
			}
			pairs[i] = pair{t: t, a: a[i]}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].t < pairs[j].t })
		var sum int64
		ok := true
		for _, p := range pairs {
			sum += p.a
			if sum > k*p.t {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
