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

	var n, p, q, r int
	if _, err := fmt.Fscan(in, &n, &p, &q, &r); err != nil {
		return
	}

	type card struct{ a, b, c int }
	cards := make([]card, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cards[i].a, &cards[i].b, &cards[i].c)
	}

	// Arrays for suffix maxima
	A := make([]int, p+2) // max b for strength >= x
	B := make([]int, p+2) // max c for strength >= x
	C := make([]int, q+2) // max c for defense >= y

	for _, crd := range cards {
		if crd.b > A[crd.a] {
			A[crd.a] = crd.b
		}
		if crd.c > B[crd.a] {
			B[crd.a] = crd.c
		}
		if crd.c > C[crd.b] {
			C[crd.b] = crd.c
		}
	}

	// suffix maxima
	for i := p - 1; i >= 1; i-- {
		if A[i+1] > A[i] {
			A[i] = A[i+1]
		}
		if B[i+1] > B[i] {
			B[i] = B[i+1]
		}
	}
	for i := q - 1; i >= 1; i-- {
		if C[i+1] > C[i] {
			C[i] = C[i+1]
		}
	}

	// Precompute prefix sums for contributions using C
	arr := make([]int64, q+1)
	for y := 1; y <= q; y++ {
		v := r - C[y]
		if v < 0 {
			v = 0
		}
		arr[y] = arr[y-1] + int64(v)
	}

	// Since C[y] is non-increasing, we can binary search positions where C[y] > val

	var ans int64
	for x := 1; x <= p; x++ {
		L := A[x] + 1
		if L > q {
			continue
		}
		// find largest y with C[y] > B[x]
		// C is non-increasing, so binary search
		idx := sort.Search(q, func(i int) bool { return C[i+1] <= B[x] })
		// idx is first index where C[idx+1] <= B[x]; so last > is idx
		last := idx // number of y with C[y] > B[x]
		if last > q {
			last = q
		}
		if last < L {
			// All y >= L have C[y] <= B[x]
			v := r - B[x]
			if v > 0 {
				ans += int64(v) * int64(q-L+1)
			}
			continue
		}
		// y in [L, last]
		if last >= L {
			ans += arr[last] - arr[L-1]
		}
		// y in [max(L, last+1), q]
		start := last + 1
		if start < L {
			start = L
		}
		if start <= q {
			v := r - B[x]
			if v > 0 {
				ans += int64(v) * int64(q-start+1)
			}
		}
	}

	fmt.Fprintln(out, ans)
}
