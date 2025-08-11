package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, d int64
	fmt.Fscan(in, &n, &d)
	a := make([]int64, n+1)
	wg := make([]int64, n+1)
	ps := make([]int64, n+1)
	for i := int64(1); i <= n; i++ {
		fmt.Fscan(in, &a[i])
		ps[i] = ps[i-1] + a[i]
	}
	for i := int64(1); i <= n; i++ {
		fmt.Fscan(in, &wg[i])
	}
	sx := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		val := ps[n] - ps[i-1] - wg[i]
		if sx[i+1] > val {
			val = sx[i+1]
		}
		sx[i] = val
	}
	z0 := sx[1]
	const inf int64 = 1 << 60
	mn := inf
	w := inf
	var z1 int64
	for i := int64(1); i < n; i++ {
		if wg[i] < mn {
			mn = wg[i]
		}
		pre := ps[i] - mn
		if pre < 0 {
			pre = 0
		}
		cand := sx[i+1] + pre
		if cand > z1 {
			z1 = cand
		}
		if i > 1 {
			v := a[i]
			if wg[i] < v {
				v = wg[i]
			}
			if v < w {
				w = v
			}
		}
	}
	cand1 := ps[n] - wg[1] - w
	if cand1 > z1 {
		z1 = cand1
	}
	z2 := ps[n] - mn
	tmp := a[n] - wg[n]
	if tmp > z2 {
		z2 = tmp
	}
	if d == 0 {
		fmt.Println(z0)
	} else if d == 1 {
		fmt.Println(z1)
	} else {
		fmt.Println(z2)
	}
}
