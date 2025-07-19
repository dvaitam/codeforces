package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	f := make([][]int, n+1)
	for i := 0; i < 2*n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		f[u] = append(f[u], v)
		f[v] = append(f[v], u)
	}

	a := make([]int, n+1)
	use := make([]bool, n+1)
	a[1] = 1

	var getRem func() bool
	getRem = func() bool {
		for i := 3; i <= n-1; i++ {
			found := false
			for _, u := range f[a[i-1]] {
				if use[u] {
					continue
				}
				ok := false
				for _, v := range f[a[i-2]] {
					if v == u {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				use[u] = true
				a[i] = u
				found = true
				break
			}
			if !found {
				return false
			}
		}
		return true
	}

	var check func() bool
	check = func() bool {
		t1, t2, t3, t4 := false, false, false, false
		for _, v := range f[a[n]] {
			if v == a[1] {
				t1 = true
			}
			if v == a[2] {
				t2 = true
			}
			if v == a[n-1] {
				t3 = true
			}
			if v == a[n-2] {
				t4 = true
			}
		}
		return t1 && t2 && t3 && t4
	}

	for p := 0; p < len(f[1]); p++ {
		for q := 0; q < len(f[1]); q++ {
			if p == q {
				continue
			}
			u := f[1][p]
			v := f[1][q]
			ok := false
			for _, x := range f[u] {
				if x == v {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
			for i := range use {
				use[i] = false
			}
			use[1], use[u], use[v] = true, true, true
			a[2] = v
			a[n] = u
			if getRem() && check() {
				for i := 1; i <= n; i++ {
					if i > 1 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, a[i])
				}
				fmt.Fprint(out, "\n")
				return
			}
		}
	}

	fmt.Fprintln(out, -1)
}
