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

const (
	refSource   = "./2122B.go"
	randomTests = 120
	totalNLimit = 180000 // within 2e5 constraint
	maxNPerCase = 50000
	maxVal      = 1_000_000_000
)

type pile struct {
	a, b, c, d int64
}

type testCase struct {
	piles []pile
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect := parseOutputs(expectRaw, len(tests))
	got := parseOutputs(gotRaw, len(tests))

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2122B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+5)

	// Deterministic small cases.
	tests = append(tests,
		testCase{piles: []pile{{1, 3, 1, 2}}},
		testCase{piles: []pile{{1, 1, 1, 2}, {3, 0, 2, 2}}},
		testCase{piles: []pile{{0, 0, 0, 0}}},
	)

	sumN := 0
	for _, tc := range tests {
		sumN += len(tc.piles)
	}

	for i := 0; i < randomTests && sumN < totalNLimit; i++ {
		remain := totalNLimit - sumN
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		sumN += n
		tests = append(tests, randomCase(n, rng))
	}

	return tests
}

func randomCase(n int, rng *rand.Rand) testCase {
	piles := make([]pile, n)

	// Distribute total zeros and ones initially.
	var totalZeros, totalOnes int64
	for i := 0; i < n; i++ {
		a := rng.Int63n(maxVal)
		b := rng.Int63n(maxVal)
		piles[i].a = a
		piles[i].b = b
		totalZeros += a
		totalOnes += b
	}

	// Create target zeros summing to totalZeros, similarly ones.
	leftZeros := totalZeros
	leftOnes := totalOnes
	for i := 0; i < n; i++ {
		if i == n-1 {
			piles[i].c = leftZeros
			piles[i].d = leftOnes
			break
		}
		zMax := leftZeros
		if zMax > maxVal {
			zMax = maxVal
		}
		targetZ := rng.Int63n(zMax + 1)
		// ensure enough zeros remain for remaining piles (at least 0 each)
		if targetZ > leftZeros-int64(n-i-1) {
			targetZ = leftZeros - int64(n-i-1)
		}
		if targetZ < 0 {
			targetZ = 0
		}
		piles[i].c = targetZ
		leftZeros -= targetZ

		oMax := leftOnes
		if oMax > maxVal {
			oMax = maxVal
		}
		targetO := rng.Int63n(oMax + 1)
		if targetO > leftOnes-int64(n-i-1) {
			targetO = leftOnes - int64(n-i-1)
		}
		if targetO < 0 {
			targetO = 0
		}
		piles[i].d = targetO
		leftOnes -= targetO
	}

	return testCase{piles: piles}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.piles)))
		sb.WriteByte('\n')
		for _, p := range tc.piles {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", p.a, p.b, p.c, p.d))
		}
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

func parseOutputs(out string, t int) []int64 {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		fail("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			fail("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
