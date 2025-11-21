package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2070-2079/2078/2078B.go"

type testCase struct {
	name  string
	input string
	n     []int
	k     []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refMaps, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refSum, err := evaluateAll(tc, refMaps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candMaps, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candSum, err := evaluateAll(tc, candMaps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if candSum != refSum {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\nexpected total distance %d, got %d\ninput:\n%s", idx+1, tc.name, refSum, candSum, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2078B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2078B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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

func parseOutput(output string, ns []int) ([][]int, error) {
	fields := strings.Fields(output)
	total := 0
	for _, v := range ns {
		total += v
	}
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d integers, got %d", total, len(fields))
	}
	res := make([][]int, len(ns))
	p := 0
	for i, n := range ns {
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[p])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[p])
			}
			res[i][j] = val
			p++
		}
	}
	return res, nil
}

func evaluateAll(tc testCase, maps [][]int) (int64, error) {
	if len(maps) != len(tc.n) {
		return 0, fmt.Errorf("expected %d testcases, got %d", len(tc.n), len(maps))
	}
	var total int64
	for i := range tc.n {
		sum, err := evaluate(tc.n[i], tc.k[i], maps[i])
		if err != nil {
			return 0, fmt.Errorf("case %d: %v", i+1, err)
		}
		total += sum
	}
	return total, nil
}

func evaluate(n int, k int64, arr []int) (int64, error) {
	if len(arr) != n {
		return 0, fmt.Errorf("expected %d values, got %d", n, len(arr))
	}
	for i, v := range arr {
		if v < 1 || v > n {
			return 0, fmt.Errorf("teleporter %d points to %d out of bounds", i+1, v)
		}
		if v == i+1 {
			return 0, fmt.Errorf("teleporter %d points to itself", i+1)
		}
	}

	// binary lifting
	maxPow := 0
	for (int64(1) << maxPow) <= k {
		maxPow++
	}
	up := make([][]int, maxPow)
	up[0] = make([]int, n)
	for i := 0; i < n; i++ {
		up[0][i] = arr[i] - 1
	}
	for p := 1; p < maxPow; p++ {
		up[p] = make([]int, n)
		for i := 0; i < n; i++ {
			up[p][i] = up[p-1][up[p-1][i]]
		}
	}

	var sum int64
	for i := 0; i < n; i++ {
		pos := i
		kk := k
		p := 0
		for kk > 0 {
			if kk&1 == 1 {
				pos = up[p][pos]
			}
			kk >>= 1
			p++
		}
		// distance from exit is n - (pos+1)
		sum += int64(n - (pos + 1))
	}
	return sum, nil
}

func buildTests() []testCase {
	tests := []testCase{
		buildFromCases("sample", []int{2, 3}, []int64{1, 2}),
		buildFromCases("min_n_k1", []int{2}, []int64{1}),
		buildFromCases("k_even_small", []int{2}, []int64{2}),
		buildFromCases("larger_k", []int{5}, []int64{15}),
		buildFromCases("two_cases", []int{4, 6}, []int64{3, 8}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildFromCases(name string, ns []int, ks []int64) testCase {
	if len(ns) != len(ks) {
		panic("ns and ks length mismatch")
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(ns)))
	sb.WriteByte('\n')
	for i := range ns {
		sb.WriteString(fmt.Sprintf("%d %d\n", ns[i], ks[i]))
	}
	return testCase{name: name, input: sb.String(), n: ns, k: ks}
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	ns := make([]int, t)
	ks := make([]int64, t)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(15) + 2
		k := rng.Int63n(1_000_000_000) + 1
		ns[i] = n
		ks[i] = k
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	}
	return testCase{name: fmt.Sprintf("random_%d", idx), input: sb.String(), n: ns, k: ks}
}
