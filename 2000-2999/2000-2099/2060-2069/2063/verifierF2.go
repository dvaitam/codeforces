package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const MOD = 998244353

type pair struct {
	l, r int
}

type testCase struct {
	n     int
	pairs []pair
}

var fact []int64
var inv_fact []int64

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func modInverse(n int64) int64 {
	return power(n, MOD-2)
}

func precompute(maxN int) {
	sz := 2*maxN + 10
	if sz < 10 {
		sz = 10
	}
	fact = make([]int64, sz)
	inv_fact = make([]int64, sz)
	fact[0] = 1
	inv_fact[0] = 1
	for i := 1; i < sz; i++ {
		fact[i] = (fact[i-1] * int64(i)) % MOD
	}
	inv_fact[sz-1] = modInverse(fact[sz-1])
	for i := sz - 2; i >= 1; i-- {
		inv_fact[i] = (inv_fact[i+1] * int64(i+1)) % MOD
	}
}

func catalan(k int) int64 {
	if k < 0 {
		return 0
	}
	res := fact[2*k]
	res = (res * inv_fact[k]) % MOD
	res = (res * inv_fact[k+1]) % MOD
	return res
}

func invCatalan(k int) int64 {
	if k < 0 {
		return 0
	}
	res := fact[k]
	res = (res * fact[k+1]) % MOD
	res = (res * inv_fact[2*k]) % MOD
	return res
}

type interval struct {
	id int
	l  int
	r  int
}

// Embedded correct solver for 2063F2
func solveReference(tests []testCase) string {
	// Determine max n for precompute
	maxN := 0
	for _, tc := range tests {
		if tc.n > maxN {
			maxN = tc.n
		}
	}
	precompute(maxN + 5)

	var out bytes.Buffer

	for ti, tc := range tests {
		n := tc.n

		arr := make([]interval, n+1)
		arr[0] = interval{0, 0, 2*n + 1}
		for i := 1; i <= n; i++ {
			arr[i] = interval{i, tc.pairs[i-1].l, tc.pairs[i-1].r}
		}

		sort.Slice(arr, func(i, j int) bool {
			if arr[i].l != arr[j].l {
				return arr[i].l < arr[j].l
			}
			return arr[i].r > arr[j].r
		})

		rVal := make([]int, n+1)
		for i := 0; i <= n; i++ {
			rVal[arr[i].id] = arr[i].r
		}

		parent := make([]int, n+1)
		st := make([]int, 0, n+1)
		st = append(st, arr[0].id)

		for i := 1; i <= n; i++ {
			for len(st) > 0 && rVal[st[len(st)-1]] < arr[i].r {
				st = st[:len(st)-1]
			}
			parent[arr[i].id] = st[len(st)-1]
			st = append(st, arr[i].id)
		}

		m := make([]int, n+1)
		for i := 0; i <= n; i++ {
			length := arr[i].r - arr[i].l + 1
			m[arr[i].id] = (length - 2) / 2
		}

		for i := 1; i <= n; i++ {
			u := arr[i].id
			p := parent[u]
			length := arr[i].r - arr[i].l + 1
			m[p] -= length / 2
		}

		ways := int64(1)
		for i := 0; i <= n; i++ {
			ways = (ways * catalan(m[i])) % MOD
		}

		ans := make([]int64, n+1)
		ans[n] = ways

		dsuParent := make([]int, n+1)
		for i := 0; i <= n; i++ {
			dsuParent[i] = i
		}

		var findSet func(int) int
		findSet = func(v int) int {
			root := v
			for root != dsuParent[root] {
				root = dsuParent[root]
			}
			curr := v
			for curr != root {
				nxt := dsuParent[curr]
				dsuParent[curr] = root
				curr = nxt
			}
			return root
		}

		for i := n; i >= 1; i-- {
			u := i
			p := findSet(parent[u])

			ways = (ways * invCatalan(m[p])) % MOD
			ways = (ways * invCatalan(m[u])) % MOD

			m[p] += m[u] + 1

			ways = (ways * catalan(m[p])) % MOD

			dsuParent[u] = p

			ans[i-1] = ways
		}

		for i := 0; i <= n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.FormatInt(ans[i], 10))
		}
		out.WriteByte('\n')

		if ti+1 < len(tests) {
			// continue
		}
	}

	return strings.TrimSpace(out.String())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	tests := generateTests()
	input := buildInput(tests)

	refOut := solveReference(tests)
	refVals, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candVals, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range refVals {
		if refVals[i] != candVals[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d token %d: expected %d got %d\n", findTest(tests, i)+1, findTokenIndex(tests, i)+1, refVals[i], candVals[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch {
	case strings.HasSuffix(path, ".go"):
		return exec.Command("go", "run", path)
	case strings.HasSuffix(path, ".py"):
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func parseOutputs(out string, tests []testCase) ([]int64, error) {
	tokens := strings.Fields(out)
	total := 0
	for _, tc := range tests {
		total += tc.n + 1
	}
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d tokens, got %d", total, len(tokens))
	}
	res := make([]int64, total)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tok)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0
	const maxSumN = 300000

	add := func(tc testCase) {
		if totalN+tc.n > maxSumN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	add(makeFromSequence("(())"))
	add(makeFromSequence("()()"))
	add(makeFromSequence("(()())"))
	add(makeFromSequence("((()))"))
	add(makeFromSequence("()(())"))

	for len(tests) < 120 && totalN < maxSumN {
		n := rng.Intn(5000) + 2
		if totalN+n > maxSumN {
			n = maxSumN - totalN
		}
		seq := randomBalanced(rng, n)
		tc := makeFromSequence(seq)
		shufflePairs(rng, tc.pairs)
		add(tc)
	}

	return tests
}

func makeFromSequence(seq string) testCase {
	stack := make([]int, 0)
	pairs := make([]pair, 0, len(seq)/2)
	for i, ch := range seq {
		pos := i + 1
		if ch == '(' {
			stack = append(stack, pos)
		} else {
			open := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pairs = append(pairs, pair{l: open, r: pos})
		}
	}
	return testCase{n: len(seq) / 2, pairs: pairs}
}

func randomBalanced(rng *rand.Rand, n int) string {
	seq := make([]byte, 0, 2*n)
	balance := 0
	for i := 0; i < 2*n; i++ {
		rem := 2*n - i
		if balance == rem {
			seq = append(seq, ')')
			balance--
			continue
		}
		if balance == 0 {
			seq = append(seq, '(')
			balance++
			continue
		}
		if rng.Intn(2) == 0 {
			seq = append(seq, '(')
			balance++
		} else {
			seq = append(seq, ')')
			balance--
		}
	}
	return string(seq)
}

func shufflePairs(rng *rand.Rand, arr []pair) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for _, p := range tc.pairs {
			fmt.Fprintf(&b, "%d %d\n", p.l, p.r)
		}
	}
	return b.String()
}

func findTest(tests []testCase, idx int) int {
	acc := 0
	for t, tc := range tests {
		lenT := tc.n + 1
		if idx < acc+lenT {
			return t
		}
		acc += lenT
	}
	return len(tests) - 1
}

func findTokenIndex(tests []testCase, idx int) int {
	acc := 0
	for _, tc := range tests {
		lenT := tc.n + 1
		if idx < acc+lenT {
			return idx - acc
		}
		acc += lenT
	}
	return 0
}
