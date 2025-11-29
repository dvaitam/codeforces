package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const solution1250GSource = `package main

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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		A := make([]int64, n+1)
		B := make([]int64, n+1)
		M := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			A[i] = A[i-1] + a[i-1]
			B[i] = B[i-1] + b[i-1]
			if A[i] < B[i] {
				M[i] = A[i]
			} else {
				M[i] = B[i]
			}
		}
		r := 0
		base := int64(0)
		resets := []int{}
		success := false
		for {
			tpos := r + 1
			for tpos <= n && B[tpos]-base < k {
				tpos++
			}
			if tpos > n {
				break
			}
			if A[tpos]-base < k {
				success = true
				break
			}
			need := A[tpos] - k + 1
			j := sort.Search(len(M[:tpos]), func(i int) bool { return M[i] >= need })
			if j >= tpos {
				break
			}
			r = j
			base = M[r]
			resets = append(resets, r)
			if r >= tpos {
				tpos = r + 1
			}
		}
		if !success {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, len(resets))
			if len(resets) > 0 {
				for i, v := range resets {
					if i > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, v)
				}
				fmt.Fprintln(out)
			}
		}
	}
}
`

// Keep the embedded reference solution reachable.
var _ = solution1250GSource

type testCase struct {
	n int
	k int64
	a []int64
	b []int64
}

const testcasesRaw = `100
9 4
2 9 7 2 9 1 5 6 6
1 8 3 9 10 1 10 10 10
7 5
8 5 10 9 4 5 3
4 4 8 4 4 5 6
4 10
6 3 2 7
4 4 10 10
4 7
3 8 5 10
2 2 5 4
1 1
6
2
8 5
5 8 7 5 4 9 6 9
10 2 6 7 6 2 9 3
10 8
10 10 7 9 6 5 9 3 3 2
3 2 5 1 4 9 6 7 9 6
2 9
10 8
1 4
9 7
1 7 7 10 9 9 4 2 2
6 7 9 5 4 8 8 3 10
10 5
9 9 8 6 2 1 6 6 8 3
9 5 1 2 8 3 8 5 9 3
3 10
2 4 9
3 8 1
2 9
6 6
7 6
4 2
7 8 10 10
9 2 10 10
10 3
7 10 5 1 9 8 7 3 2 5
4 9 1 3 4 1 6 1 4 7
8 3
9 1 10 6 3 7 10 8
4 9 3 8 5 6 7 8
9 6
7 9 9 9 9 8 1 7 3
6 4 8 5 2 7 3 5 6
1 7
10
10 3
8 10 6 3 8 1 1 4 8 10
4 3 4 6 7 6 1 3 6 1
3 10
7 10 9
5 5 5
5 8
8 7 8 4 10
1 9 5 9 8
3 10
5 5 3
1 3 5
10 5
7 10 8 4 1 8 9 1 3 6
2 3 10 1 6 1 10 8 8 4
1 7
10
9 10
5 4 8 6 10 2 10 10 3
4 10 5 8 8 8 9 1 4
8 10
7 10 9 10 9 6 8 3
5 5 7 1 2 7 5 8
2 5
10 8
9 10
10 10 8 8 4 10 8 3 10 3
4 7 2 6 1 9 6 4 3 8
3 5
10 5 8
7 9 10
9 9
7 8 8 4 3 7 1 1 9
7 9 1 10 6 6 4 4 1
6 4
10 1 9 1 8 1
3 3 2 1 2 8
10 3
9 1 7 6 7 9 6 5 8 10
6 1 4 9 5 7 3 5 3 5
3 4
4 9 1
8 6 9
9 8
4 4 2 9 8 6 7 3 4
10 5 10 1 9 1 9 5 8
9 2
8 7 9 3 5 8 3 5 6
3 9 10 1 8 5 10 8 5
5 6
4 10 7 2 2
6 2 2 10 5
3 6
4 6 9
9 1 7
6 2
8 10 6 2 1 9
3 6 5 8 2 6
9 7
9 6 2 7 1 9 8 6 3
5 6 8 8 8 5 8 2 5
8 9
9 10 7 3 3 2 1 5
3 10 4 8 7 5 7 2
7 6
4 6 5 4 4 4 1
1 10 6 10 7 2 8
6 4
5 5 2 10 2 10
1 5 4 1 1 9
3 2
3 9 1
9 7 10
10 1
4 5 8 1 8 5 10 6 1 1
4 4 1 8 3 9 3 9 7 3
3 3
6 1 9
5 9 9
2 4
6 1
7 2
1 5
7
4
8 2
2 7 10 8 3 3 2 8
2 1 10 10 7 10 1 1
10 10
7 8 6 8 3 9 2 6 3 5
9 3 4 10 3 8 2 10 7 2
3 4
5 1 8
2 6 7
4 6
9 7 2 9
10 4 1 9
1 10
4
2 2
10 8
5 9
6 3 9 6 3
1 7 10 4 3
8 9
6 3 8 1 3 10 3 1
1 8 7 8 5 5 4 9
10 1
1 2 2 4 3 4 6 6 7 7
9 7 10 1 4 10 10 1 1 9
4 5
10 6 1 3
1 2 5 4
8 9
5 5 7 7 7 5 1 7
7 3 10 1 5 4 3 4
5 6
10 1 3 10 6
6 1 2 8 5
3 7
10 10 8
10 2 6
8 10
2 3 7 4 4 8 2 8
9 6 10 6 10 2 9 8
10 2
8 6 8 10 8 7 1 4 9 2
9 9 9 5 3 4 3 4 9 10
6 2
3 5 6 1 4 4
7 10 5 9 7 10
6 1
1 1 7 2 3 6
10 3 8 10 4 8
4 9
3 9 2 9
9 4 1 8
1 2
9
4 4
9 7 4 1
1 1 9 6
4 10
4 2 6 3
4 7 3 6
6 8
3 1 2 9 9 7
6 5 10 8 7 9
8 1
3 6 6 7 1 9 9 10
7 1 1 5 6 9 10 7
4 6
6 5 4 1
6 10 6 4
6 2
9 3 8 5 5 1
8 1 8 1 3 9
2 4
8 9
1 7 4 3 9 10 6 2
3 8 7 3 6 7 3 9
6 3
9 1 2 1 10 9
3 3 7 7 8 6
7 9
8 3 7 7 1 8 7
8 7 3 2 2 2 1
2 5
6 8
3 6
10 9 9 7 3 10
10 10 10 2 2 10
9 10
1 7 6 9 8 9 8 4 9
4 6 4 3 7 6 9 5 7
5 4
9 2 10 6 6
1 3 10 5 7
6 6
4 9 8 6 5 10
8 5 10 1 8 7
9 5
10 7 4 8 1 9 8 1 3
3 10 5 8 3 5 4 6 1
5 8
4 6 6 1 8
8 3 8 6 9
5 6
8 5 2 6 9
7 9 7 3 8
10 6
9 7 1 2 2 2 3 9 8 7
7 8 8 7 2 6 10 4 9 7
3 1
4 8 1
5 1 8
2 7
9 1
7 2
4 5
2 9
4 8
5 10
9 7 9 4 4
6 10 9 2 6
9 10
10 1 4 9 6 2 10 9 4 1
4 4 2 5 4 10 4 7 8 10
10 4
5 4 6 5 1 8 7 7 6 2
9 6 2 7 3 8 5 4 1 4
4 7
1 4 1 10
9 6 10 6
10 1
4 3 6 10 5 3 8 6 8 10
6 6 7 9 7 10 3 10 9 1
9 1
1 2 7 8 1 2 9 10 6
7 10 7 2 6 2 6 2 6
10 2
2 3 1 8 7 2 1 7 7 3
4 6 4 6 9 4 5 9 9 9
2 10
5 9
8 2
4 3 4 5 7 9 10 10
6 8 5 8 2 1 8 10
9 9
10 8 10 1 10 5 6 9 6
9 6 8 6 8 2 7 3 6
9 3
9 10 5 4 4 10 3 9 4
9 3 6 10 3 6 7 1 8
8 7
6 5 6 7 1 9 9 7
6 8 1 9 1 8 6 5
2 1
2 6
7 8
8 9 8 8 9 10 6 8
6 6 10 9 2 7 5 6
10 1
6 10 3 10 6 6 8 10 9 10
5 8 3 6 8 8 9 2 10 9
1 3
1
8 4
3 9 8 2 9 2 1 4
5 3 7 7 7 6 7 5
6 2
8 8 5 7 6 8
4 1 5 5 3 3
5 4
8 10 7 7 8
1 1 6 7 3
4 4
3 4 5 1
8 6 9 1
10 2
4 10 6 5 5 4 9 3 7 8
8 10 2 9 4 2 1 10 9 3
8 1
4 4 8 8 9 3 8 5
7 6 1 8 7 4 5 8
4 10
9 9 2 9
6 10 2 1
3 4
1 2 8
1 2 6
7 6
6 6 2 6 10 4 5
5 6 4 1 2 1 1
8 10
1 2 3 5 7 6 4 8
3 10 7 7 5 1 5 4
5 8
2 7 10 6 3
4 10 7 2 10
2 4
10 5
6 2
10 8 3 4 8 6
2 1 4 2 2 10
9 6
6 9 7 9 7 1 10 3 10
9 4 4 6 4 4 9 10 10
4 7
7 7 6 8
6 10 10 5
8 8
6 9 1 9 1 4 6 1
6 3 2 4 3 9 7 1
9 7
8 7 5 4 8 8 4 3 3
10 4 5 5 2 7 6 9 3
10 5
6 9 3 6 8 8 7 4 8 9
2 8 3 6 7 4 8 8 7 10
1 4
5
2 7
1 5
3 3
8 7 1
2 5 7
7 8
2 6 8 7 3 7 2 10
3 8 4 9 3 9 10 9
4 6
9 9 2 1
1 5 6 4
6 4
9 2 6 8 2 5
10 8 7 9 6 9
10 1
3 3 4 6 3 9 9 6 3 3
5 7 3 1 1 5 8 5 1 5
8 3
9 7 3 9 5 4 9 4
3 1 1 8 10 7 5 3
6 9
6 6 9 3 7 3
1 1 2 1 7 3
10 1
10 5 9 5 4 6 8 4 4 10
2 1 5 6 7 2 1 7 5 3
10 10
7 7 6 9 3 3 9 2 2 2
1 5 8 5 10 9 1 10 6 10
8 9
6 1 9 9 7 1 1 10
2 5 1 10 2 3 6 7
10 1
4 3 1 6 10 10 6 3 10 8
1 6 4 3 6 4 6 1 7 6
7 1
10 9 2 9 9 10 10
10 9 4 6 8 9 6
10 7
10 1 3 9 3 9 3 9 8 10
1 10 4 3 5 1 6 5 9 2
8 7
1 5 5 5 7 6 8 2
6 6 5 5 6 2 4 10
10 4
4 2 7 7 1 7 5 4 4 8
3 6 7 8 1 7 10 4 6 9
5 5
4 3 3 9 2
8 10 10 10 3
4 3
10 7 6 1
8 6 3 6
9 9
9 2 9 10 10 5 7 5 1
7 5 8 7 5 8 6 8 7
6 8
6 5 2 3 7 6
7 4 6 2 6 2
1 6
6
7 8
2 9 2 10 9 2 4 10
3 4 1 4 8 8 6 5
3 8
4 4 3
9 3 5
1 3
8
1 8
7
7 5
4 2
`

func parseTestcases() []testCase {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		panic(err)
	}
	res := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
		}
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &b[j])
		}
		res[i] = testCase{n: n, k: k, a: a, b: b}
	}
	return res
}

func solveExpected(tc testCase) string {
	n := tc.n
	k := tc.k
	a := tc.a
	b := tc.b

	A := make([]int64, n+1)
	B := make([]int64, n+1)
	M := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		A[i] = A[i-1] + a[i-1]
		B[i] = B[i-1] + b[i-1]
		if A[i] < B[i] {
			M[i] = A[i]
		} else {
			M[i] = B[i]
		}
	}
	r := 0
	base := int64(0)
	var resets []int
	success := false
	for {
		tpos := r + 1
		for tpos <= n && B[tpos]-base < k {
			tpos++
		}
		if tpos > n {
			break
		}
		if A[tpos]-base < k {
			success = true
			break
		}
		need := A[tpos] - k + 1
		j := sort.Search(len(M[:tpos]), func(i int) bool { return M[i] >= need })
		if j >= tpos {
			break
		}
		r = j
		base = M[r]
		resets = append(resets, r)
		if r >= tpos {
			tpos = r + 1
		}
	}
	if !success {
		return "-1"
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, len(resets))
	if len(resets) > 0 {
		for i, v := range resets {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return strings.TrimSpace(sb.String())
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.k)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		input := buildInput(tc)
		expect := solveExpected(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\ninput:\n%s", i+1, err, string(out), input)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
