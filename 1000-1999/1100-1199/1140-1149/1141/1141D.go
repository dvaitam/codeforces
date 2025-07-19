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
	var a, b string
	fmt.Fscan(in, &a, &b)

	// buckets for 'a'-'z' (0-25) and '?' (26)
	A := make([][]int, 27)
	B := make([][]int, 27)
	for i := 0; i < n; i++ {
		ai := 26
		if a[i] != '?' {
			ai = int(a[i] - 'a')
		}
		bi := 26
		if b[i] != '?' {
			bi = int(b[i] - 'a')
		}
		A[ai] = append(A[ai], i+1)
		B[bi] = append(B[bi], i+1)
	}

	var ans [][2]int
	// match same letters and wildcards
	for c := 0; c < 26; c++ {
		// same char
		for len(A[c]) > 0 && len(B[c]) > 0 {
			ai := A[c][len(A[c])-1]
			A[c] = A[c][:len(A[c])-1]
			bi := B[c][len(B[c])-1]
			B[c] = B[c][:len(B[c])-1]
			ans = append(ans, [2]int{ai, bi})
		}
		// A char with B '?'
		for len(A[c]) > 0 && len(B[26]) > 0 {
			ai := A[c][len(A[c])-1]
			A[c] = A[c][:len(A[c])-1]
			bi := B[26][len(B[26])-1]
			B[26] = B[26][:len(B[26])-1]
			ans = append(ans, [2]int{ai, bi})
		}
		// A '?' with B char
		for len(A[26]) > 0 && len(B[c]) > 0 {
			ai := A[26][len(A[26])-1]
			A[26] = A[26][:len(A[26])-1]
			bi := B[c][len(B[c])-1]
			B[c] = B[c][:len(B[c])-1]
			ans = append(ans, [2]int{ai, bi})
		}
	}
	// '?' with '?'
	for len(A[26]) > 0 && len(B[26]) > 0 {
		ai := A[26][len(A[26])-1]
		A[26] = A[26][:len(A[26])-1]
		bi := B[26][len(B[26])-1]
		B[26] = B[26][:len(B[26])-1]
		ans = append(ans, [2]int{ai, bi})
	}

	fmt.Fprintln(out, len(ans))
	for _, p := range ans {
		fmt.Fprintln(out, p[0], p[1])
	}
}
