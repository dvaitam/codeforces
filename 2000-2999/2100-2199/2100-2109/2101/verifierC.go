package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	randomTests   = 120
	maxTotalN     = 200000
	maxCaseLength = 60000
)

type testCase struct {
	n int
	a []int
}

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	for i := range p {
		p[i] = i
	}
	return &dsu{p}
}

func (d *dsu) find(x int) int {
	if d.parent[x] == x {
		return x
	}
	d.parent[x] = d.find(d.parent[x])
	return d.parent[x]
}

func (d *dsu) union(x, y int) {
	rx := d.find(x)
	ry := d.find(y)
	if rx != ry {
		d.parent[rx] = ry
	}
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// refSolve is the correct embedded reference solver for 2101C.
// It processes multiple test cases from a single input string.
func refSolve(input string) []int64 {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(reader, &t)

	results := make([]int64, 0, t)

	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)

		type Item struct {
			id       int
			deadline int
			weight   int
		}

		items := make([]Item, n)
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(reader, &a)
			items[i] = Item{
				id:       i + 1,
				deadline: 2 * a,
				weight:   iabs(2*(i+1) - n - 1),
			}
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].weight == items[j].weight {
				return items[i].id > items[j].id
			}
			return items[i].weight > items[j].weight
		})

		dsuObj := newDSU(2 * n)
		accepted := make([]int, 0, n)

		for _, item := range items {
			d := item.deadline
			if d > 2*n {
				d = 2 * n
			}
			slot := dsuObj.find(d)
			if slot > 0 {
				accepted = append(accepted, item.id)
				dsuObj.union(slot, slot-1)
			}
		}

		sort.Ints(accepted)
		m := len(accepted)
		k := m / 2

		var profit int64
		for i := 0; i < k; i++ {
			profit -= int64(accepted[i])
		}
		for i := m - k; i < m; i++ {
			profit += int64(accepted[i])
		}

		results = append(results, profit)
	}

	return results
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	expect := refSolve(input)

	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}
	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v\n%s", err, gotRaw)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			tc := tests[i]
			fail("mismatch on test %d: expected %d, got %d (n=%d, a[:min(10)]%v)", i+1, expect[i], got[i], tc.n, preview(tc.a, 10))
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+5)

	totalN := 0
	add := func(a []int) {
		n := len(a)
		if n == 0 || totalN+n > maxTotalN {
			return
		}
		tests = append(tests, testCase{n: n, a: a})
		totalN += n
	}

	// Deterministic coverage including samples and edge cases.
	add([]int{1, 2, 1, 2})
	add([]int{2, 2})
	add([]int{1, 2, 1, 5, 1, 2, 2, 1, 1, 2})
	add([]int{1, 5, 2, 8, 4, 1, 4, 2})
	add([]int{1})
	add([]int{1, 1, 1, 1})
	add([]int{4, 3, 2, 1})

	for len(tests) < randomTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxLen := maxCaseLength
		if maxLen > remain {
			maxLen = remain
		}
		if maxLen < 1 {
			break
		}
		n := rng.Intn(maxLen) + 1

		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1 // within [1, n]
		}
		add(a)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func preview(a []int, k int) []int {
	if k > len(a) {
		k = len(a)
	}
	cpy := make([]int, k)
	copy(cpy, a[:k])
	return cpy
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
