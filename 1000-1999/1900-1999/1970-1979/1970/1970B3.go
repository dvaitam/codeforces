package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]struct{ fs, sn int }, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i].fs)
		a[i].sn = i
	}
	sort.Slice(a[1:], func(i, j int) bool {
		return a[i+1].fs < a[j+1].fs
	})
	x := 0
	for i := 1; i < n; i++ {
		if a[i].fs == a[i+1].fs {
			x = i
			break
		}
	}
	ansX := make([]int, n+1)
	ansY := make([]int, n+1)
	anss := make([]int, n+1)
	v := make([]int, n+1)
	if x > 0 || a[1].fs == 0 {
		if x > 0 {
			a[1], a[x] = a[x], a[1]
			a[2], a[x+1] = a[x+1], a[2]
		}
		// first element
		ansX[a[1].sn] = 1
		ansY[a[1].sn] = 1
		v[1] = 1
		if a[1].fs == 0 {
			anss[a[1].sn] = a[1].sn
		} else {
			anss[a[1].sn] = a[2].sn
		}
		// others
		for i := 2; i <= n; i++ {
			idx := a[i].sn
			if a[i].fs == 0 {
				ansX[idx] = i
				ansY[idx] = 1
				v[i] = 1
				anss[idx] = a[i].sn
			} else if a[i].fs < i {
				ansX[idx] = i
				y := v[i-a[i].fs]
				ansY[idx] = y
				v[i] = y
				anss[idx] = a[i-a[i].fs].sn
			} else {
				ansX[idx] = i
				y := a[i].fs - i + 2
				ansY[idx] = y
				v[i] = y
				anss[idx] = a[1].sn
			}
		}
	} else {
		if n == 2 {
			fmt.Fprintln(out, "NO")
			return
		}
		// special case
		ansX[a[n].sn] = 1
		ansY[a[n].sn] = 1
		anss[a[n].sn] = a[n-1].sn
		ansX[a[n-1].sn] = n
		ansY[a[n-1].sn] = 2
		anss[a[n-1].sn] = a[1].sn
		for i := 1; i <= n-2; i++ {
			idx := a[i].sn
			ansX[idx] = i + 1
			ansY[idx] = 1
			anss[idx] = a[n].sn
		}
	}
	fmt.Fprintln(out, "YES")
	for i := 1; i <= n; i++ {
		fmt.Fprintf(out, "%d %d\n", ansX[i], ansY[i])
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(out, "%d", anss[i])
	}
	out.WriteByte('\n')
}
