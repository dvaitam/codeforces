package main

import (
	"bufio"
	"fmt"
	"os"
)

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), size: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func sieve(n int) []int {
	spf := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i*i <= n {
				for j := i * i; j <= n; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	return spf
}

func factorize(x int, spf []int) []int {
	res := []int{}
	for x > 1 {
		p := spf[x]
		res = append(res, p)
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func uniqueInts(arr []int) []int {
	m := make(map[int]struct{}, len(arr))
	out := make([]int, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	arr := make([]int, n)
	maxVal := 1000001
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i]+1 > maxVal {
			maxVal = arr[i] + 1
		}
	}
	spf := sieve(maxVal)
	d := newDSU(maxVal + 1)

	// union primes for each value
	for _, x := range arr {
		fac := uniqueInts(factorize(x, spf))
		if len(fac) == 0 {
			continue
		}
		base := fac[0]
		for _, p := range fac[1:] {
			d.union(base, p)
		}
	}

	roots := make([]int, n)
	plusRoots := make([][]int, n)
	edges := make(map[[2]int]struct{})

	for idx, x := range arr {
		fac := uniqueInts(factorize(x, spf))
		baseRoot := d.find(fac[0])
		roots[idx] = baseRoot

		// gather all DSU roots from x and x+1
		unionSet := []int{baseRoot}
		for _, p := range fac[1:] {
			rp := d.find(p)
			if rp != baseRoot {
				unionSet = append(unionSet, rp)
			}
		}
		fac2 := uniqueInts(factorize(x+1, spf))
		tmp := make([]int, 0, len(fac2))
		for _, p := range fac2 {
			rp := d.find(p)
			tmp = append(tmp, rp)
			unionSet = append(unionSet, rp)
		}
		plusRoots[idx] = uniqueInts(tmp)
		unionSet = uniqueInts(unionSet)
		for i := 0; i < len(unionSet); i++ {
			for j := i + 1; j < len(unionSet); j++ {
				a, b := unionSet[i], unionSet[j]
				if a > b {
					a, b = b, a
				}
				edges[[2]int{a, b}] = struct{}{}
			}
		}
	}

	for ; q > 0; q-- {
		var s, t int
		fmt.Fscan(in, &s, &t)
		s--
		t--
		if roots[s] == roots[t] {
			fmt.Fprintln(out, 0)
			continue
		}
		checkSetS := append([]int{roots[s]}, plusRoots[s]...)
		checkSetT := append([]int{roots[t]}, plusRoots[t]...)
		found := false
		for _, rs := range checkSetS {
			for _, rt := range checkSetT {
				if rs == rt {
					found = true
					break
				}
				a, b := rs, rt
				if a > b {
					a, b = b, a
				}
				if _, ok := edges[[2]int{a, b}]; ok {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
