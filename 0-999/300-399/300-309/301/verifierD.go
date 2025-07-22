package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestCaseD struct {
	input    string
	expected string
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// --- Reference solution logic copied from 301D.go ---

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *BIT) Sum(i int) int {
	if i <= 0 {
		return 0
	}
	s := 0
	for ; i > 0; i -= i & -i {
		s += b.tree[i]
	}
	return s
}

type Query struct {
	y     int
	id    int
	coeff int
}

func solveLocal(n, m int, p []int, L, R []int) []int64 {
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pos[p[i-1]] = i
	}
	edgesByU := make([][]int, n+1)
	for v := 1; v <= n; v++ {
		u := pos[v]
		for k := v + v; k <= n; k += v {
			edgesByU[u] = append(edgesByU[u], pos[k])
		}
	}
	subsByX := make([][]Query, n+1)
	for i := 0; i < m; i++ {
		l, r := L[i], R[i]
		subsByX[r] = append(subsByX[r], Query{y: r, id: i, coeff: 1})
		if l > 0 {
			subsByX[l-1] = append(subsByX[l-1], Query{y: r, id: i, coeff: -1})
			subsByX[r] = append(subsByX[r], Query{y: l - 1, id: i, coeff: -1})
			subsByX[l-1] = append(subsByX[l-1], Query{y: l - 1, id: i, coeff: 1})
		}
	}
	bit := NewBIT(n)
	ansEdges := make([]int64, m)
	for x := 0; x <= n; x++ {
		if x > 0 {
			for _, y := range edgesByU[x] {
				bit.Add(y, 1)
			}
		}
		for _, q := range subsByX[x] {
			cnt := bit.Sum(q.y)
			ansEdges[q.id] += int64(q.coeff) * int64(cnt)
		}
	}
	res := make([]int64, m)
	for i := 0; i < m; i++ {
		res[i] = ansEdges[i] + int64(R[i]-L[i]+1)
	}
	return res
}

func genTests() []TestCaseD {
	rand.Seed(3)
	tests := make([]TestCaseD, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(50) + 1
		m := rand.Intn(50) + 1
		perm := rand.Perm(n)
		for i := range perm {
			perm[i] += 1
		}
		L := make([]int, m)
		R := make([]int, m)
		for i := 0; i < m; i++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			L[i], R[i] = l, r
		}
		ans := solveLocal(n, m, perm, L, R)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d %d\n", L[i], R[i])
		}
		outLines := make([]string, m)
		for i, v := range ans {
			outLines[i] = strconv.FormatInt(v, 10)
		}
		tests = append(tests, TestCaseD{input: sb.String(), expected: strings.Join(outLines, "\n")})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	passed := 0
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			continue
		}
		cleaned := strings.TrimSpace(out)
		if cleaned != tc.expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, tc.expected, cleaned)
		} else {
			passed++
		}
	}
	fmt.Printf("passed %d/%d tests\n", passed, len(tests))
	if passed != len(tests) {
		os.Exit(1)
	}
}
