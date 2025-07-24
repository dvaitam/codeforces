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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	cand := make([]int, n)
	for i := range cand {
		cand[i] = i + 1
	}
	contains := func(set map[int]struct{}, v int) bool {
		_, ok := set[v]
		return ok
	}

	for len(cand) > 2 {
		m := len(cand)
		third := m / 3
		if third == 0 {
			third = 1
		}
		A := cand[:third]
		B := cand[third:min(2*third, m)]
		C := cand[min(2*third, m):]

		set1 := make([]int, 0, len(A)+len(B))
		set1 = append(set1, A...)
		set1 = append(set1, B...)

		set2 := make([]int, 0, len(A)+len(C))
		set2 = append(set2, A...)
		set2 = append(set2, C...)

		// ask first set
		fmt.Fprint(out, "? ", len(set1))
		for _, v := range set1 {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
		out.Flush()
		var r1 string
		if _, err := fmt.Fscan(in, &r1); err != nil {
			return
		}

		// ask second set
		fmt.Fprint(out, "? ", len(set2))
		for _, v := range set2 {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
		out.Flush()
		var r2 string
		if _, err := fmt.Fscan(in, &r2); err != nil {
			return
		}

		setMap1 := make(map[int]struct{}, len(set1))
		for _, v := range set1 {
			setMap1[v] = struct{}{}
		}
		setMap2 := make(map[int]struct{}, len(set2))
		for _, v := range set2 {
			setMap2[v] = struct{}{}
		}

		var next []int
		for _, x := range cand {
			ok1 := (r1 == "YES" && contains(setMap1, x)) || (r1 == "NO" && !contains(setMap1, x))
			ok2 := (r2 == "YES" && contains(setMap2, x)) || (r2 == "NO" && !contains(setMap2, x))
			if ok1 || ok2 {
				next = append(next, x)
			}
		}
		cand = next
	}

	for _, v := range cand {
		fmt.Fprintln(out, "!", v)
		out.Flush()
		var r string
		if _, err := fmt.Fscan(in, &r); err != nil {
			return
		}
		if r == ":)" {
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
