package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type query struct {
	typ int
	l   int
	r   int
}

type test struct {
	input  string
	output string
}

func solve(n int, v []int64, qs []query) string {
	orig := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		orig[i] = orig[i-1] + v[i-1]
	}
	u := make([]int64, n)
	copy(u, v)
	sort.Slice(u, func(i, j int) bool { return u[i] < u[j] })
	sorted := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		sorted[i] = sorted[i-1] + u[i-1]
	}
	var sb strings.Builder
	for i, q := range qs {
		var ans int64
		if q.typ == 1 {
			ans = orig[q.r] - orig[q.l-1]
		} else {
			ans = sorted[q.r] - sorted[q.l-1]
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d", ans)
	}
	return sb.String()
}

func generateTests() []test {
	rand.Seed(2)
	var tests []test
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		v := make([]int64, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			v[j] = int64(rand.Intn(100) + 1)
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v[j])
		}
		sb.WriteByte('\n')
		m := rand.Intn(20) + 1
		fmt.Fprintf(&sb, "%d\n", m)
		qs := make([]query, m)
		for j := 0; j < m; j++ {
			typ := rand.Intn(2) + 1
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			qs[j] = query{typ, l, r}
			fmt.Fprintf(&sb, "%d %d %d\n", typ, l, r)
		}
		out := solve(n, v, qs)
		tests = append(tests, test{sb.String(), out})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.output {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
