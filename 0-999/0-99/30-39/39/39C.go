package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	l, r, length int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var N int
	if _, err := fmt.Fscan(in, &N); err != nil {
		return
	}
	A := make([]Node, N+1)
	for i := 1; i <= N; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		A[i].l = x - y + 1
		A[i].r = x + y
		A[i].length = A[i].r - A[i].l + 1
	}
	P := make([]int, N+1)
	Q := make([]int, N+1)
	for i := 1; i <= N; i++ {
		P[i] = i
		Q[i] = i
	}
	// sort by length
	sort.Slice(P[1:], func(i, j int) bool {
		return A[P[i+1]].length < A[P[j+1]].length
	})
	// sort by left end
	sort.Slice(Q[1:], func(i, j int) bool {
		return A[Q[i+1]].l < A[Q[j+1]].l
	})
	const INF = 2000000005
	f := make([]int, N+2)
	link := make([]int, N+2)
	g := make([]int, N+1)
	res := make([][]int, N+1)
	// phase 1: compute g and res
	for k := 1; k <= N; k++ {
		i := P[k]
		// reset f and link
		for j := 0; j <= N; j++ {
			f[j] = INF
			link[j] = 0
		}
		f[0] = A[i].l - 1
		p := 0
		for h := 1; h <= N; h++ {
			j := Q[h]
			if A[j].r > A[i].r || A[j].l < A[i].l {
				continue
			}
			for p < N && f[p+1] < A[j].l {
				p++
			}
			idx := p + g[j]
			if idx <= N && A[j].r < f[idx] {
				f[idx] = A[j].r
				link[idx] = j
			}
		}
		// find max j
		jmax := N
		for jmax > 0 && f[jmax] == INF {
			jmax--
		}
		// build res for i
		seq := make([]int, 0, jmax+1)
		for h := jmax; h > 0; h -= g[link[h]] {
			seq = append(seq, link[h])
		}
		// append i as root
		seq = append(seq, i)
		res[i] = seq
		g[i] = jmax + 1
	}
	// phase 2: final DP
	for j := 0; j <= N; j++ {
		f[j] = INF
		link[j] = 0
	}
	f[0] = A[1].l - 1
	p := 0
	for h := 1; h <= N; h++ {
		j := Q[h]
		for p < N && f[p+1] < A[j].l {
			p++
		}
		idx := p + g[j]
		if idx <= N && A[j].r < f[idx] {
			f[idx] = A[j].r
			link[idx] = j
		}
	}
	// find jmax
	J := N
	for J > 0 && f[J] == INF {
		J--
	}
	// reconstruct
	var Out []int
	var Get func(x int)
	Get = func(x int) {
		Out = append(Out, x)
		seq := res[x]
		// children are seq[0:len-1]
		for _, ch := range seq[:len(seq)-1] {
			Get(ch)
		}
	}
	for h := J; h > 0; h -= g[link[h]] {
		Get(link[h])
	}
	// sort and output
	sort.Ints(Out)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, J)
	for i, v := range Out {
		if i+1 < len(Out) {
			fmt.Fprintf(w, "%d ", v)
		} else {
			fmt.Fprintf(w, "%d\n", v)
		}
	}
}
