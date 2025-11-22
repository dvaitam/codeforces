package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxK = 60                  // 2^60 > 1e18
	inf  = uint64(1<<62) - 1e6 // large sentinel to avoid overflow
)

type test struct {
	input, expected string
}

type caseData struct {
	a []uint64
	b []uint64
}

func ceilDiv(a, b uint64) uint64 {
	if a >= inf {
		return inf
	}
	return (a + b - 1) / b
}

func value(arr []uint64, idx int) uint64 {
	if idx < len(arr) {
		return arr[idx]
	}
	return 1
}

func ensureTailOne(arr []uint64) []uint64 {
	if len(arr) == 0 || arr[len(arr)-1] != 1 {
		arr = append(arr, 1)
	}
	return arr
}

func trimTrailingOnes(arr []uint64) []uint64 {
	for len(arr) > 1 && arr[len(arr)-1] == 1 && arr[len(arr)-2] == 1 {
		arr = arr[:len(arr)-1]
	}
	return arr
}

func effectiveLen(arr []uint64) int {
	if len(arr) == 0 {
		return 0
	}
	if arr[len(arr)-1] == 1 {
		return len(arr) - 1
	}
	return len(arr)
}

// res[m] = min_{i+j=m} max(a[i], b[j]); a and b are non-increasing.
func combine(a, b []uint64) []uint64 {
	effA := effectiveLen(a)
	effB := effectiveLen(b)
	limit := effA + effB
	if limit > maxK {
		limit = maxK
	}
	res := make([]uint64, limit+1)
	p := 0
	for m := 0; m <= limit; m++ {
		if p > m {
			p = m
		}
		best := value(a, p)
		vb := value(b, m-p)
		if vb > best {
			best = vb
		}
		for p < m {
			cand := value(a, p+1)
			v2 := value(b, m-p-1)
			if v2 > cand {
				cand = v2
			}
			if cand <= best {
				p++
				best = cand
			} else {
				break
			}
		}
		res[m] = best
	}
	res = ensureTailOne(res)
	res = trimTrailingOnes(res)
	return res
}

func buildPows(base uint64) [maxK + 1]uint64 {
	var p [maxK + 1]uint64
	p[0] = 1
	for i := 1; i <= maxK; i++ {
		if p[i-1] >= inf/base {
			p[i] = inf
		} else {
			p[i] = p[i-1] * base
			if p[i] > inf {
				p[i] = inf
			}
		}
	}
	return p
}

func computeNeed(aVal, bVal uint64, child []uint64) []uint64 {
	pows := buildPows(bVal)
	need := make([]uint64, maxK+1)
	for i := range need {
		need[i] = inf
	}

	maxStart := aVal
	c0 := value(child, 0)
	if c0 > maxStart {
		maxStart = c0
	}

	childEff := effectiveLen(child)
	for t := 0; t <= maxK; t++ {
		factor := pows[t]
		if factor >= maxStart {
			for k := t; k <= maxK; k++ {
				need[k] = 1
			}
			break
		}
		limitRem := maxK - t
		if childEff < limitRem {
			limitRem = childEff
		}
		for rem := 0; rem <= limitRem; rem++ {
			cur := value(child, rem)
			if aVal > cur {
				cur = aVal
			}
			req := ceilDiv(cur, factor)
			k := t + rem
			if req < need[k] {
				need[k] = req
			}
			if req == 1 {
				for kk := k; kk <= maxK && need[kk] > 1; kk++ {
					need[kk] = 1
				}
				break
			}
		}
	}

	for k := 1; k <= maxK; k++ {
		if need[k] > need[k-1] {
			need[k] = need[k-1]
		}
	}
	need = ensureTailOne(need)
	need = trimTrailingOnes(need)
	return need
}

func buildCartesian(b []uint64) (int, []int, []int) {
	n := len(b)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && b[stack[len(stack)-1]] > b[i] {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			parent[i] = stack[len(stack)-1]
		}
		if last != -1 {
			parent[last] = i
		}
		stack = append(stack, i)
	}

	left := make([]int, n)
	right := make([]int, n)
	for i := 0; i < n; i++ {
		left[i], right[i] = -1, -1
	}
	root := -1
	for i, p := range parent {
		if p == -1 {
			root = i
			continue
		}
		if i < p {
			left[p] = i
		} else {
			right[p] = i
		}
	}
	return root, left, right
}

func solveOne(a, b []uint64) int {
	root, left, right := buildCartesian(b)
	order := make([]int, 0, len(a))
	stack := []int{root}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		if left[v] != -1 {
			stack = append(stack, left[v])
		}
		if right[v] != -1 {
			stack = append(stack, right[v])
		}
	}

	needs := make([][]uint64, len(a))
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		ln := []uint64{1}
		if left[v] != -1 {
			ln = needs[left[v]]
		}
		rn := []uint64{1}
		if right[v] != -1 {
			rn = needs[right[v]]
		}
		child := combine(ln, rn)
		needs[v] = computeNeed(a[v], b[v], child)
	}

	ans := 0
	rootNeed := needs[root]
	for ans < len(rootNeed) && rootNeed[ans] > 1 {
		ans++
	}
	return ans
}

func formatInput(cases []caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		n := len(cs.a)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range cs.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func expectedOutput(cases []caseData) string {
	var sb strings.Builder
	for i, cs := range cases {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", solveOne(cs.a, cs.b)))
	}
	return sb.String()
}

func buildTest(cases []caseData) test {
	return test{
		input:    formatInput(cases),
		expected: expectedOutput(cases),
	}
}

func randVal(rng *rand.Rand, hi uint64) uint64 {
	return uint64(rng.Int63n(int64(hi))) + 1
}

func genRandomCase(rng *rand.Rand, n int) caseData {
	a := make([]uint64, n)
	b := make([]uint64, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			a[i] = randVal(rng, 1_000_000_000_000_000_000)
		} else {
			a[i] = randVal(rng, 1_000_000)
		}
		if rng.Intn(6) == 0 {
			b[i] = 1_000_000_000_000_000_000 - uint64(rng.Intn(1000))
		} else {
			b[i] = uint64(rng.Intn(1_000_000-1) + 2) // ensure at least 2
		}
		if b[i] < 2 {
			b[i] = 2
		}
	}
	return caseData{a: a, b: b}
}

func generateTests() []test {
	var tests []test
	// small hand-crafted cases
	tests = append(tests, buildTest([]caseData{{a: []uint64{1}, b: []uint64{2}}}))
	tests = append(tests, buildTest([]caseData{{a: []uint64{5}, b: []uint64{2}}}))
	tests = append(tests, buildTest([]caseData{{a: []uint64{5, 4, 2}, b: []uint64{3, 2, 2}}}))
	tests = append(tests, buildTest([]caseData{
		{a: []uint64{7, 7, 7, 7}, b: []uint64{9, 9, 9, 9}},
		{a: []uint64{10, 1, 10}, b: []uint64{3, 2, 4}},
	}))

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 40; i++ {
		caseCnt := rng.Intn(4) + 1
		totalN := 0
		cases := make([]caseData, caseCnt)
		for j := 0; j < caseCnt; j++ {
			n := rng.Intn(25) + 1
			totalN += n
			if totalN > 180 { // keep inputs quick while still varied
				n = 1 + rng.Intn(10)
			}
			cases[j] = genRandomCase(rng, n)
		}
		tests = append(tests, buildTest(cases))
	}

	// a moderately larger case to cover heavier paths
	large := genRandomCase(rand.New(rand.NewSource(123)), 200)
	tests = append(tests, buildTest([]caseData{large}))
	return tests
}

func prepareBinary(path, tag string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	bin := filepath.Join(os.TempDir(), tag)
	cmd := exec.Command("go", "build", "-o", bin, path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", func() {}, fmt.Errorf("build failed: %v\n%s", err, out)
	}
	cleanup := func() { os.Remove(bin) }
	return bin, cleanup, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(8*time.Second, func() {
		cmd.Process.Kill()
	})
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := prepareBinary(os.Args[1], "cand2048F")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(cand, tc.input)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := strings.TrimSpace(tc.expected)
		if exp != got {
			fmt.Printf("Wrong answer on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
