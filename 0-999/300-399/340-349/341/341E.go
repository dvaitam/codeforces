package main

import (
	"bufio"
	"fmt"
	"os"
)

// PII represents a pair of integers: fi for value, se for original index
type PII struct{ fi, se int }

var ret []PII

// add performs a move: from box b to box a, recording the move,
// then updates values: a.fi doubles, b.fi decreases by previous a.fi
func add(a, b *PII) {
	ret = append(ret, PII{a.se, b.se})
	b.fi -= a.fi
	a.fi <<= 1
}

// sort3 sorts a length-3 array of PII by fi asc, se asc
func sort3(t *[3]PII) {
	if t[0].fi > t[1].fi || (t[0].fi == t[1].fi && t[0].se > t[1].se) {
		t[0], t[1] = t[1], t[0]
	}
	if t[1].fi > t[2].fi || (t[1].fi == t[2].fi && t[1].se > t[2].se) {
		t[1], t[2] = t[2], t[1]
	}
	if t[0].fi > t[1].fi || (t[0].fi == t[1].fi && t[0].se > t[1].se) {
		t[0], t[1] = t[1], t[0]
	}
}

// work recursively balances three boxes until smallest is zero
func work(t *[3]PII) {
	if t[0].fi == 0 {
		return
	}
	x := t[1].fi / t[0].fi
	for x > 0 {
		if x&1 == 1 {
			add(&t[0], &t[1])
		} else {
			add(&t[0], &t[2])
		}
		x >>= 1
	}
	sort3(t)
	work(t)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var Q []PII
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x > 0 {
			Q = append(Q, PII{x, i})
		}
	}
	num := len(Q)
	if num <= 1 {
		fmt.Println(-1)
		return
	}
	if num == 2 {
		fmt.Println(0)
		return
	}
	head := 0
	for i := 0; i < num-2; i++ {
		t := [3]PII{Q[head], Q[head+1], Q[head+2]}
		head += 3
		sort3(&t)
		work(&t)
		Q = append(Q, t[1], t[2])
	}
	// output moves
	fmt.Println(len(ret))
	for _, p := range ret {
		fmt.Println(p.fi, p.se)
	}
}
