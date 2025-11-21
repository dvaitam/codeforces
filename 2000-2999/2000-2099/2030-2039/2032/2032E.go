package main

import (
	"bufio"
	"fmt"
	"os"
)

func solveCase(a []int64) ([]int64, bool) {
	n := len(a)
	if n == 1 {
		return []int64{0}, true
	}
	p := make([]int64, n)
	q := make([]int64, n)
	r := make([]int64, n)
	q[1] = 1
	for i := 1; i <= n-2; i++ {
		p[i+1] = 1 - p[i-1] - 2*p[i]
		q[i+1] = -q[i-1] - 2*q[i]
		r[i+1] = -a[i] - r[i-1] - 2*r[i]
	}

	A1 := p[n-1] - 1
	B1 := q[n-1] + 1
	C1 := r[n-1] + a[0]

	A2 := p[n-2] + 2*p[n-1] - 1
	B2 := q[n-2] + 2*q[n-1]
	C2 := r[n-2] + 2*r[n-1] + a[n-1]

	D := A1*B2 - A2*B1
	if D == 0 {
		return nil, false
	}

	Fnum := -C1*B2 + C2*B1
	Xnum := -A1*C2 + A2*C1
	if Fnum%D != 0 || Xnum%D != 0 {
		return nil, false
	}
	F := Fnum / D
	x := Xnum / D

	v := make([]int64, n)
	v[0] = 0
	v[1] = x
	for i := 1; i <= n-2; i++ {
		v[i+1] = F - a[i] - v[i-1] - 2*v[i]
	}

	minVal := v[0]
	for _, val := range v {
		if val < minVal {
			minVal = val
		}
	}
	if minVal < 0 {
		shift := -minVal
		for i := range v {
			v[i] += shift
		}
	}

	return v, true
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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		v, ok := solveCase(a)
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		for i, val := range v {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val)
		}
		fmt.Fprintln(out)
	}
}
