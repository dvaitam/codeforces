package main

import (
	"bufio"
	"fmt"
	"os"
)

type testCase struct {
	n int
	m int
	s []int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	cases := make([]testCase, t)
	maxN := 0
	maxVal := 0
	for i := 0; i < t; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		s := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &s[j])
			if s[j] > maxVal {
				maxVal = s[j]
			}
		}
		if n > maxN {
			maxN = n
		}
		cases[i] = testCase{n: n, m: m, s: s}
	}

	divIdx := make([][]int, maxN+1)
	for d := 1; d <= maxN; d++ {
		for mult := d * 2; mult <= maxN; mult += d {
			divIdx[mult] = append(divIdx[mult], d)
		}
	}

	divVal := make([][]int, maxVal+1)
	for d := 1; d <= maxVal; d++ {
		for mult := d; mult <= maxVal; mult += d {
			divVal[mult] = append(divVal[mult], d)
		}
	}

	for _, tc := range cases {
		if tc.m == 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		desc := make([]int, tc.m)
		for i := 0; i < tc.m; i++ {
			desc[i] = tc.s[tc.m-1-i]
		}

		ans := make([]int, tc.n+1)
		forbid := make([]bool, maxVal+1)
		used := make([]int, 0)

		ok := true
		for i := 1; i <= tc.n && ok; i++ {
			used = used[:0]
			for _, d := range divIdx[i] {
				val := ans[d]
				if val == 0 {
					continue
				}
				if !forbid[val] {
					forbid[val] = true
					used = append(used, val)
				}
			}

			found := false
			for _, val := range desc {
				isValid := true
				for _, dv := range divVal[val] {
					if forbid[dv] {
						isValid = false
						break
					}
				}
				if isValid {
					ans[i] = val
					found = true
					break
				}
			}
			for _, v := range used {
				forbid[v] = false
			}
			if !found {
				ok = false
			}
		}

		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
