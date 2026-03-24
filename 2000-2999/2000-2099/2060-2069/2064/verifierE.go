package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---- Embedded correct solver for 2064 E ----

func solveReference(input string) string {
	data := []byte(input)
	pos := 0
	maxLen := len(data)
	readNext := func() byte {
		if pos >= maxLen {
			return 0
		}
		b := data[pos]
		pos++
		return b
	}

	readInt := func() int {
		b := readNext()
		for b != 0 && (b < '0' || b > '9') {
			b = readNext()
		}
		if b == 0 {
			return 0
		}
		res := 0
		for b >= '0' && b <= '9' {
			res = res*10 + int(b-'0')
			b = readNext()
		}
		return res
	}

	const mod int64 = 998244353
	const maxN = 200005

	parent := make([]int, maxN)
	sz := make([]int, maxN)
	color := make([]int, maxN)
	prev := make([]int, maxN)
	nxt := make([]int, maxN)
	x := make([]int, maxN)
	p := make([]int, maxN)
	c := make([]int, maxN)

	var out bytes.Buffer

	t := readInt()
	for t > 0 {
		t--
		n := readInt()
		for i := 1; i <= n; i++ {
			p[i] = readInt()
		}
		for i := 1; i <= n; i++ {
			c[i] = readInt()
		}

		for i := 1; i <= n; i++ {
			x[p[i]] = i
		}

		for i := 1; i <= n; i++ {
			parent[i] = i
			sz[i] = 1
			color[i] = c[i]
			prev[i] = i - 1
			nxt[i] = i + 1
		}
		prev[1] = 0
		nxt[n] = 0

		find := func(i int) int {
			root := i
			for parent[root] != root {
				root = parent[root]
			}
			curr := i
			for curr != root {
				nNode := parent[curr]
				parent[curr] = root
				curr = nNode
			}
			return root
		}

		merge := func(u, v int) {
			parent[v] = u
			sz[u] += sz[v]
			nxt[u] = nxt[v]
			if nxt[v] != 0 {
				prev[nxt[v]] = u
			}
		}

		curr := 1
		for nxt[curr] != 0 {
			nx := nxt[curr]
			if color[curr] == color[nx] {
				merge(curr, nx)
			} else {
				curr = nx
			}
		}

		ans := int64(1)

		for k := 1; k <= n; k++ {
			idx := x[k]
			b := find(idx)

			ans = (ans * int64(sz[b])) % mod
			sz[b]--

			if sz[b] == 0 {
				pBlk := prev[b]
				nBlk := nxt[b]

				if pBlk != 0 {
					nxt[pBlk] = nBlk
				}
				if nBlk != 0 {
					prev[nBlk] = pBlk
				}

				if pBlk != 0 && nBlk != 0 && color[pBlk] == color[nBlk] {
					merge(pBlk, nBlk)
				}
			}
		}

		out.WriteString(strconv.FormatInt(ans, 10))
		out.WriteByte('\n')
	}
	return out.String()
}

// ---- Verifier infrastructure ----

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()

	for idx, tc := range tests {
		refOut := solveReference(tc.input)
		refAns, err := parseAnswers(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if !equalSlices(refAns, candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sreference:\n%s\ncandidate:\n%s", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample", "4\n1\n1\n1\n5\n3 4 1 2 5\n1 1 1 1 1\n5\n4 2 3 1 5\n2 1 4 1 5\n40\n29 15 20 35 37 31 27 1 32 36 38 25 22 8 16 7 3 28 11 12 23 4 14 9 39 13 10 30 6 2 24 17 19 5 34 18 33 26 40 21\n3 1 2 2 1 2 3 1 1 1 1 2 1 3 1 1 3 1 1 1 2 2 1 3 3 3 2 3 2 2 2 2 1 3 2 1 1 2 2 2\n"),
		buildSingleCase("n1", []int{1}, []int{1}),
		buildSingleCase("sorted colors same", []int{1, 2, 3, 4}, []int{7, 7, 7, 7}),
		buildSingleCase("permutation reversed", []int{4, 3, 2, 1}, []int{1, 2, 3, 4}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 25; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSingleCase(name string, p, c []int) testCase {
	if len(p) != len(c) {
		panic("p and c length mismatch")
	}
	n := len(p)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String(), answers: 1}
}

func newTestCase(name, input string) testCase {
	q, err := countAnswers(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, answers: q}
}

func countAnswers(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read t: %v", err)
	}
	if t <= 0 {
		return 0, fmt.Errorf("non-positive t: %d", t)
	}
	return t, nil
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(25) + 1
		p := rng.Perm(n)
		for j := range p {
			p[j]++
		}
		c := make([]int, n)
		for j := range c {
			c[j] = rng.Intn(n) + 1
		}
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j, v := range p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for j, v := range c {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	name := fmt.Sprintf("random_%d", idx)
	return testCase{name: name, input: sb.String(), answers: t}
}
