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

const refSource = "./2032D.go"

type testCase struct {
	name  string
	input string
}

type treeCase struct {
	n       int
	parents []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier2032D-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	tests := []testCase{
		buildTest("basic", []treeCase{
			caseFromLengths([]int{2, 1}),
			caseFromLengths([]int{3, 1, 1}),
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomTest(fmt.Sprintf("random_small_%d", i+1), rng, 5, 200))
	}
	for i := 0; i < 3; i++ {
		tests = append(tests, randomTest(fmt.Sprintf("random_mid_%d", i+1), rng, 10, 2000))
	}
	tests = append(tests, randomTest("random_large", rng, 20, 10000))

	return tests
}

func buildTest(name string, cases []treeCase) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, cs := range cases {
		fmt.Fprintln(&b, cs.n)
		for i, val := range cs.parents {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, val)
		}
		fmt.Fprintln(&b)
	}
	return testCase{name: name, input: b.String()}
}

func randomTest(name string, rng *rand.Rand, maxTests int, maxTotalN int) testCase {
	if maxTests < 1 {
		maxTests = 1
	}
	t := rng.Intn(maxTests) + 1
	var cases []treeCase
	sumN := 0
	for len(cases) < t && sumN < maxTotalN {
		remaining := maxTotalN - sumN
		if remaining < 4 {
			break
		}
		maxN := remaining
		if maxN > 2000 {
			maxN = 2000
		}
		n := rng.Intn(maxN-3) + 4
		parents := randomParents(rng, n-1)
		cases = append(cases, treeCase{n: n, parents: parents})
		sumN += n
	}
	if len(cases) == 0 {
		cases = append(cases, caseFromLengths([]int{2, 1}))
	}
	return buildTest(name, cases)
}

func caseFromLengths(lengths []int) treeCase {
	parents := parentsFromLengths(lengths)
	return treeCase{n: len(parents) + 1, parents: parents}
}

func parentsFromLengths(lengths []int) []int {
	total := 0
	for _, v := range lengths {
		total += v
	}
	parent := make([]int, total+1) // 1-indexed
	remaining := make([]int, len(lengths))
	pending := make(map[int]int)

	idx := 1
	for chainID, length := range lengths {
		if length <= 0 {
			panic("invalid chain length")
		}
		parent[idx] = 0
		remaining[chainID] = length - 1
		if remaining[chainID] > 0 {
			pending[idx] = chainID
		}
		idx++
	}

	for parentIdx := 1; parentIdx <= total; parentIdx++ {
		if idx > total {
			break
		}
		chainID, ok := pending[parentIdx]
		if !ok {
			continue
		}
		delete(pending, parentIdx)
		parent[idx] = parentIdx
		remaining[chainID]--
		if remaining[chainID] > 0 {
			pending[idx] = chainID
		}
		idx++
	}

	res := make([]int, total)
	for i := 1; i <= total; i++ {
		res[i-1] = parent[i]
	}
	return res
}

func randomParents(rng *rand.Rand, total int) []int {
	if total < 3 {
		total = 3
	}
	for {
		maxChains := total
		if maxChains > 8 {
			maxChains = 8
		}
		if maxChains < 1 {
			maxChains = 1
		}
		k := rng.Intn(maxChains) + 1
		minSum := 2 + (k - 1)
		if minSum > total {
			continue
		}
		lengths := make([]int, k)
		lengths[0] = 2
		for i := 1; i < k; i++ {
			lengths[i] = 1
		}
		remain := total - minSum
		for remain > 0 {
			idx := rng.Intn(k)
			lengths[idx]++
			remain--
		}
		return parentsFromLengths(lengths)
	}
}
