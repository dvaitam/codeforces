package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func dfs(u int, vals []int, left []int, right []int, isLeaf []bool) ([]int, int) {
	if isLeaf[u] {
		return []int{vals[u]}, 0
	}
	a, hl := dfs(left[u], vals, left, right, isLeaf)
	b, hr := dfs(right[u], vals, left, right, isLeaf)
	v := append(a, b...)
	sort.Ints(v)
	h := hl
	if hr > h {
		h = hr
	}
	h++
	if h%2 == 1 {
		v = v[1:]
	} else {
		v = v[:len(v)-1]
	}
	return v, h
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var t int
	fmt.Fscan(rdr, &t)
	var outputs []string
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(rdr, &n)
		vals := make([]int, n)
		left := make([]int, n)
		right := make([]int, n)
		isLeaf := make([]bool, n)
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(rdr, &a)
			if a >= 0 {
				vals[i] = a
				isLeaf[i] = true
			} else {
				var l, r int
				fmt.Fscan(rdr, &l, &r)
				left[i] = l
				right[i] = r
			}
		}
		res, _ := dfs(0, vals, left, right, isLeaf)
		if len(res) > 0 {
			outputs = append(outputs, fmt.Sprintf("%d", res[0]))
		} else {
			outputs = append(outputs, "0")
		}
	}
	return strings.Join(outputs, "\n")
}

func genTreeCase() string {
	n := rand.Intn(5) + 1
	vals := make([]int, n)
	left := make([]int, n)
	right := make([]int, n)
	isLeaf := make([]bool, n)
	for i := 0; i < n; i++ {
		if i > 0 && rand.Intn(2) == 0 {
			// leaf
			vals[i] = rand.Intn(10)
			isLeaf[i] = true
		} else if i > 0 {
			// non-leaf, connect to previous nodes less than i
			l := rand.Intn(i)
			r := rand.Intn(i)
			left[i] = l
			right[i] = r
		} else {
			// root may be leaf or not
			vals[i] = rand.Intn(10)
			isLeaf[i] = true
		}
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, 1)
	fmt.Fprintln(&buf, n)
	for i := 0; i < n; i++ {
		if isLeaf[i] {
			fmt.Fprintf(&buf, "%d\n", vals[i])
		} else {
			fmt.Fprintf(&buf, "-1 %d %d\n", left[i], right[i])
		}
	}
	return buf.String()
}

func generateCases() []testCase {
	rand.Seed(6)
	cases := []testCase{}
	fixed := []string{
		"1\n1\n5\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		inp := genTreeCase()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
