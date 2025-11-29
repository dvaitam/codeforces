package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1139ESource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var N, M int
	fmt.Fscan(reader, &N, &M)
	P := make([]int, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &P[i])
	}
	C := make([]int, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &C[i])
	}
	var D int
	fmt.Fscan(reader, &D)
	K := make([]int, N+1)
	found := make([]bool, N+1)
	for i := 1; i <= D; i++ {
		fmt.Fscan(reader, &K[i])
		found[K[i]] = true
	}
	newD := D
	for i := 1; i <= N; i++ {
		if !found[i] {
			newD++
			K[newD] = i
		}
	}
	G := make([][]int, M+1)
	cnt := make([]bool, (M+1)*(M+1))
	used := make([]bool, M+1)
	Le := make([]int, M+1)
	Ri := make([]int, M+1)
	for i := range Le {
		Le[i] = -1
	}
	var pairUp func(node int) bool
	pairUp = func(node int) bool {
		if used[node] {
			return false
		}
		used[node] = true
		for _, to := range G[node] {
			if Le[to] == -1 {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		for _, to := range G[node] {
			if pairUp(Le[to]) {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		return false
	}
	ans := make([]int, D+1)
	L := -1
	for i := newD; i >= 1; i-- {
		if i <= D {
			ans[i] = L + 1
		}
		p := P[K[i]]
		c := C[K[i]]
		idx := c*(M+1) + p
		if !cnt[idx] {
			G[p] = append(G[p], c)
		}
		for ok := true; ok; {
			ok = false
			for j := range used {
				used[j] = false
			}
			if pairUp(L + 1) {
				L++
				ok = true
			}
		}
		cnt[idx] = true
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= D; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1139ESource

type testCase struct {
	n int
	m int
	p []int
	c []int
	d int
	k []int
}

var testcases = []testCase{
	{n: 2, m: 2, p: []int{1, 2}, c: []int{3, 1}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{0}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{0, 5, 2}, d: 2, k: []int{1, 3}},
	{n: 3, m: 2, p: []int{2, 1, 2}, c: []int{5, 3, 4}, d: 1, k: []int{1}},
	{n: 2, m: 2, p: []int{2, 1}, c: []int{4, 2}, d: 1, k: []int{2}},
	{n: 5, m: 3, p: []int{3, 1, 2, 2, 3}, c: []int{2, 3, 3, 1, 1}, d: 3, k: []int{3, 1, 4}},
	{n: 1, m: 1, p: []int{1}, c: []int{4}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{1, 5}, d: 2, k: []int{2, 1}},
	{n: 3, m: 2, p: []int{2, 1, 2}, c: []int{0, 0, 5}, d: 1, k: []int{2}},
	{n: 5, m: 5, p: []int{2, 1, 3, 2, 3}, c: []int{3, 0, 0, 2, 5}, d: 1, k: []int{3}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{1, 5, 3}, d: 3, k: []int{3, 1, 2}},
	{n: 5, m: 2, p: []int{2, 2, 1, 2, 2}, c: []int{4, 1, 2, 4, 0}, d: 3, k: []int{1, 4, 5}},
	{n: 3, m: 2, p: []int{2, 1, 2}, c: []int{1, 3, 1}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{4}, d: 1, k: []int{1}},
	{n: 5, m: 4, p: []int{2, 3, 1, 1, 3}, c: []int{3, 5, 1, 3, 1}, d: 2, k: []int{4, 5}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{5, 3, 1, 3}, d: 2, k: []int{1, 4}},
	{n: 3, m: 2, p: []int{1, 1, 1}, c: []int{3, 2, 1}, d: 2, k: []int{1, 2}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{3, 1, 4, 1, 2}, d: 3, k: []int{4, 3, 2}},
	{n: 5, m: 2, p: []int{2, 2, 2, 2, 1}, c: []int{2, 5, 5, 1, 0}, d: 4, k: []int{5, 2, 4, 1}},
	{n: 2, m: 2, p: []int{2, 2}, c: []int{3, 1}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{3, 3}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{4, 3}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{4}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{4}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{1, 2, 3}, d: 1, k: []int{2}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{3, 5, 3, 0}, d: 2, k: []int{3, 2}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{5, 1, 2, 3}, d: 3, k: []int{3, 4, 1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{4, 0}, d: 2, k: []int{2, 1}},
	{n: 5, m: 2, p: []int{2, 2, 2, 1, 1}, c: []int{5, 0, 5, 0, 4}, d: 3, k: []int{5, 4, 3}},
	{n: 3, m: 3, p: []int{2, 3, 3}, c: []int{3, 3, 1}, d: 2, k: []int{3, 2}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{5, 4, 1}, d: 2, k: []int{2, 1}},
	{n: 5, m: 3, p: []int{3, 3, 1, 1, 2}, c: []int{5, 1, 2, 5, 2}, d: 3, k: []int{2, 3, 1}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{4, 3, 5, 4, 0}, d: 1, k: []int{2}},
	{n: 4, m: 2, p: []int{2, 2, 1, 2}, c: []int{5, 1, 2, 2}, d: 2, k: []int{3, 4}},
	{n: 5, m: 4, p: []int{4, 4, 3, 3, 4}, c: []int{3, 4, 0, 2, 1}, d: 5, k: []int{4, 3, 2, 5, 1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{4, 2}, d: 1, k: []int{2}},
	{n: 4, m: 2, p: []int{2, 1, 1, 1}, c: []int{0, 2, 1, 0}, d: 2, k: []int{4, 3}},
	{n: 5, m: 5, p: []int{5, 3, 1, 5, 2}, c: []int{1, 0, 1, 0, 2}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{1, 5, 3}, d: 3, k: []int{3, 1, 2}},
	{n: 1, m: 1, p: []int{1}, c: []int{1}, d: 1, k: []int{1}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{4, 3, 4, 5}, d: 3, k: []int{4, 2, 1}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{2, 1, 0, 1}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{5, 4}, d: 1, k: []int{2}},
	{n: 4, m: 3, p: []int{1, 3, 3, 1}, c: []int{0, 3, 4, 2}, d: 3, k: []int{4, 3, 1}},
	{n: 2, m: 2, p: []int{2, 2}, c: []int{5, 2}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{2, 0, 1}, d: 3, k: []int{1, 2, 3}},
	{n: 1, m: 1, p: []int{1}, c: []int{2}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{1}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{5}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{5, 4, 3}, d: 1, k: []int{1}},
	{n: 4, m: 4, p: []int{1, 3, 3, 2}, c: []int{5, 5, 3, 2}, d: 2, k: []int{2, 1}},
	{n: 4, m: 4, p: []int{3, 3, 4, 4}, c: []int{2, 4, 1, 2}, d: 3, k: []int{1, 4, 2}},
	{n: 3, m: 3, p: []int{2, 3, 3}, c: []int{1, 4, 2}, d: 1, k: []int{1}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{5, 2, 0, 3, 4}, d: 2, k: []int{2, 1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{2, 4}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{2, 0}, d: 2, k: []int{2, 1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{0, 1}, d: 1, k: []int{2}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{4, 4, 3, 4}, d: 3, k: []int{3, 2, 4}},
	{n: 1, m: 1, p: []int{1}, c: []int{1}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{3, 2, 4}, d: 3, k: []int{1, 3, 2}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{0, 1, 2}, d: 3, k: []int{3, 2, 1}},
	{n: 1, m: 1, p: []int{1}, c: []int{1}, d: 1, k: []int{1}},
	{n: 5, m: 4, p: []int{3, 3, 1, 4, 2}, c: []int{3, 0, 5, 5, 5}, d: 4, k: []int{2, 1, 5, 3}},
	{n: 5, m: 5, p: []int{1, 1, 4, 2, 5}, c: []int{0, 1, 4, 1, 4}, d: 3, k: []int{5, 1, 3}},
	{n: 1, m: 1, p: []int{1}, c: []int{1}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{3, 5}, d: 2, k: []int{2, 1}},
	{n: 3, m: 3, p: []int{3, 3, 3}, c: []int{2, 2, 4}, d: 2, k: []int{3, 1}},
	{n: 5, m: 2, p: []int{2, 1, 2, 2, 1}, c: []int{0, 5, 4, 1, 1}, d: 2, k: []int{5, 2}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{2, 4, 0, 2, 3}, d: 3, k: []int{2, 4, 5}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{2, 4, 4, 2, 0}, d: 2, k: []int{4, 1}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{1, 5, 5, 0}, d: 1, k: []int{1}},
	{n: 1, m: 1, p: []int{1}, c: []int{0}, d: 1, k: []int{1}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{0, 1}, d: 2, k: []int{1, 2}},
	{n: 5, m: 3, p: []int{3, 2, 3, 2, 1}, c: []int{4, 4, 0, 4, 5}, d: 5, k: []int{1, 5, 4, 3, 2}},
	{n: 5, m: 4, p: []int{3, 3, 2, 3, 4}, c: []int{1, 4, 3, 1, 3}, d: 5, k: []int{5, 1, 3, 4, 2}},
	{n: 5, m: 5, p: []int{5, 4, 3, 2, 3}, c: []int{1, 2, 2, 2, 5}, d: 2, k: []int{5, 4}},
	{n: 5, m: 3, p: []int{2, 2, 3, 1, 2}, c: []int{0, 1, 2, 5, 1}, d: 2, k: []int{4, 1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{5, 5, 4}, d: 1, k: []int{1}},
	{n: 2, m: 2, p: []int{2, 2}, c: []int{1, 5}, d: 2, k: []int{2, 1}},
	{n: 5, m: 1, p: []int{1, 1, 1, 1, 1}, c: []int{3, 0, 3, 1, 3}, d: 1, k: []int{5}},
	{n: 2, m: 2, p: []int{2, 2}, c: []int{5, 1}, d: 1, k: []int{2}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{2, 4, 3}, d: 3, k: []int{3, 1, 2}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{4, 1}, d: 1, k: []int{1}},
	{n: 3, m: 1, p: []int{1, 1, 1}, c: []int{3, 5, 2}, d: 2, k: []int{2, 3}},
	{n: 3, m: 2, p: []int{1, 2, 2}, c: []int{4, 0, 5}, d: 2, k: []int{2, 3}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{0, 2}, d: 1, k: []int{1}},
	{n: 5, m: 2, p: []int{2, 1, 2, 2, 2}, c: []int{3, 4, 5, 0, 3}, d: 1, k: []int{4}},
	{n: 3, m: 2, p: []int{2, 1, 1}, c: []int{2, 2, 1}, d: 2, k: []int{3, 2}},
	{n: 5, m: 2, p: []int{2, 2, 2, 2, 1}, c: []int{5, 5, 3, 2, 4}, d: 5, k: []int{2, 5, 1, 4, 3}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{2, 1}, d: 2, k: []int{2, 1}},
	{n: 5, m: 2, p: []int{2, 1, 1, 1, 2}, c: []int{0, 5, 5, 5, 1}, d: 2, k: []int{4, 1}},
	{n: 5, m: 5, p: []int{2, 2, 2, 1, 2}, c: []int{4, 1, 1, 0, 4}, d: 4, k: []int{3, 5, 2, 4}},
	{n: 5, m: 2, p: []int{2, 1, 2, 2, 1}, c: []int{0, 2, 5, 0, 3}, d: 1, k: []int{2}},
	{n: 2, m: 1, p: []int{1, 1}, c: []int{2, 0}, d: 1, k: []int{1}},
	{n: 4, m: 4, p: []int{4, 1, 2, 2}, c: []int{0, 1, 5, 3}, d: 4, k: []int{4, 2, 3, 1}},
	{n: 3, m: 2, p: []int{1, 2, 2}, c: []int{5, 4, 2}, d: 3, k: []int{2, 1, 3}},
	{n: 5, m: 2, p: []int{1, 2, 1, 1, 1}, c: []int{4, 0, 4, 5, 0}, d: 3, k: []int{1, 5, 3}},
	{n: 4, m: 1, p: []int{1, 1, 1, 1}, c: []int{1, 0, 2, 4}, d: 4, k: []int{3, 4, 2, 1}},
	{n: 1, m: 1, p: []int{1}, c: []int{4}, d: 1, k: []int{1}},
	{n: 5, m: 3, p: []int{3, 3, 1, 2, 1}, c: []int{1, 1, 2, 5, 0}, d: 2, k: []int{2, 4}},
}

func solveCase(tc testCase) []int {
	n := tc.n
	m := tc.m
	mm := m
	P := make([]int, n+1)
	C := make([]int, n+1)
	for i := 1; i <= n; i++ {
		P[i] = tc.p[i-1]
		C[i] = tc.c[i-1]
		if P[i] > mm {
			mm = P[i]
		}
		if C[i] > mm {
			mm = C[i]
		}
	}
	D := tc.d
	K := make([]int, n+1)
	found := make([]bool, n+1)
	for i := 1; i <= D; i++ {
		K[i] = tc.k[i-1]
		found[K[i]] = true
	}
	newD := D
	for i := 1; i <= n; i++ {
		if !found[i] {
			newD++
			K[newD] = i
		}
	}
	G := make([][]int, mm+1)
	cnt := make([]bool, (mm+1)*(mm+1))
	used := make([]bool, mm+1)
	Le := make([]int, mm+1)
	Ri := make([]int, mm+1)
	for i := range Le {
		Le[i] = -1
	}
	var pairUp func(node int) bool
	pairUp = func(node int) bool {
		if used[node] {
			return false
		}
		used[node] = true
		for _, to := range G[node] {
			if Le[to] == -1 {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		for _, to := range G[node] {
			if pairUp(Le[to]) {
				Le[to] = node
				Ri[node] = to
				return true
			}
		}
		return false
	}

	ans := make([]int, D+1)
	L := -1
	for i := newD; i >= 1; i-- {
		if i <= D {
			ans[i] = L + 1
		}
		p := P[K[i]]
		c := C[K[i]]
		idx := c*(mm+1) + p
		if !cnt[idx] {
			G[p] = append(G[p], c)
		}
		for ok := true; ok; {
			ok = false
			for j := range used {
				used[j] = false
			}
			if pairUp(L + 1) {
				L++
				ok = true
			}
		}
		cnt[idx] = true
	}
	return ans[1:]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d\n%s\n%s\n%d\n%s\n",
			tc.n, tc.m,
			intSliceToString(tc.p),
			intSliceToString(tc.c),
			tc.d,
			intSliceToString(tc.k),
		)
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		gotLines := strings.Fields(string(out))
		if len(gotLines) != len(want) {
			fmt.Printf("test %d failed: expected %d lines got %d\n", idx+1, len(want), len(gotLines))
			os.Exit(1)
		}
		for i, g := range gotLines {
			if g != fmt.Sprintf("%d", want[i]) {
				fmt.Printf("test %d failed at line %d: expected %d got %s\n", idx+1, i+1, want[i], g)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

func intSliceToString(arr []int) string {
	parts := make([]string, len(arr))
	for i, v := range arr {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ")
}
