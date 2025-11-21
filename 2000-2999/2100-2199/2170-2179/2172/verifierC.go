package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2170-2179/2172/2172C.go"

type testCase struct {
	k     int64
	n     int
	l     int
	radii []int64
	name  string
}

type fraction struct {
	num *big.Int
	den *big.Int
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := prepareBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for i, tc := range tests {
		input := singleInput(tc)
		refOut, err := runSolution(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
		expect, err := parseOutput(refOut, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}

		candOut, err := runSolution(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\nInput:\n%s", i+1, tc.name, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\nInput:\n%sOutput:\n%s", i+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalFrac(expect[0], got[0]) {
			fmt.Fprintf(os.Stderr, "Test %d (%s) failed.\nExpected: %s\nGot: %s\nInput:\n%s", i+1, tc.name, formatFrac(expect[0]), formatFrac(got[0]), input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2172C-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runSolution(binaryPath, input string) (string, error) {
	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(output string, expected int) ([]fraction, error) {
	tokens := strings.Fields(strings.TrimSpace(output))
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	res := make([]fraction, expected)
	for i, tok := range tokens {
		frac, err := parseFraction(tok)
		if err != nil {
			return nil, fmt.Errorf("token %d (%q) invalid: %v", i+1, tok, err)
		}
		res[i] = frac
	}
	return res, nil
}

func parseFraction(token string) (fraction, error) {
	if token == "" {
		return fraction{}, fmt.Errorf("empty token")
	}
	if strings.Contains(token, "/") {
		parts := strings.Split(token, "/")
		if len(parts) != 2 {
			return fraction{}, fmt.Errorf("invalid fraction format")
		}
		num := new(big.Int)
		den := new(big.Int)
		if _, ok := num.SetString(parts[0], 10); !ok {
			return fraction{}, fmt.Errorf("invalid numerator")
		}
		if _, ok := den.SetString(parts[1], 10); !ok {
			return fraction{}, fmt.Errorf("invalid denominator")
		}
		if den.Sign() == 0 {
			return fraction{}, fmt.Errorf("zero denominator")
		}
		if den.Sign() < 0 {
			num.Neg(num)
			den.Neg(den)
		}
		return fraction{num: num, den: den}, nil
	}
	num := new(big.Int)
	if _, ok := num.SetString(token, 10); !ok {
		return fraction{}, fmt.Errorf("invalid integer")
	}
	return fraction{num: num, den: big.NewInt(1)}, nil
}

func equalFrac(a, b fraction) bool {
	left := new(big.Int).Mul(a.num, b.den)
	right := new(big.Int).Mul(b.num, a.den)
	return left.Cmp(right) == 0
}

func formatFrac(f fraction) string {
	if f.den.Cmp(big.NewInt(1)) == 0 {
		return f.num.String()
	}
	return fmt.Sprintf("%s/%s", f.num.String(), f.den.String())
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 30)
	tests = append(tests,
		testCase{name: "sample1", k: 15, n: 4, l: 3, radii: []int64{7, 5, 3, 1}},
		testCase{name: "sample2", k: 14, n: 6, l: 1, radii: []int64{7, 5, 4, 3, 2, 1}},
		testCase{name: "sample3", k: 14, n: 2, l: 2, radii: []int64{5, 4}},
		testCase{name: "sample4", k: 22, n: 4, l: 4, radii: []int64{4, 4, 1, 1}},
		testCase{name: "sample5", k: 13, n: 3, l: 3, radii: []int64{6, 3, 1}},
	)

	tests = append(tests,
		testCase{name: "all_equal", k: 50, n: 5, l: 2, radii: []int64{5, 5, 5, 5, 5}},
		testCase{name: "large_k", k: 1_000_000_000, n: 5, l: 1, radii: []int64{1_000_000_000, 900_000_000, 800_000_000, 2, 1}},
		testCase{name: "many_enclosed", k: 100, n: 8, l: 3, radii: []int64{50, 45, 40, 35, 30, 25, 20, 15}},
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 25
	for len(tests) < total {
		tests = append(tests, randomTest(rng, len(tests)+1))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(180) + 2
	l := rng.Intn(min(n, 200)) + 1
	k := rng.Int63n(1_000_000_000) + 1
	radii := make([]int64, n)
	for i := range radii {
		radii[i] = rng.Int63n(1_000_000_000) + 1
	}
	sort.Slice(radii, func(i, j int) bool { return radii[i] > radii[j] })
	return testCase{k: k, n: n, l: l, radii: radii, name: fmt.Sprintf("rand_%d", idx)}
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.k, tc.n, tc.l)
	writeArray(&sb, tc.radii)
	return sb.String()
}

func writeArray(sb *strings.Builder, arr []int64) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
