package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource2062F = "2062F.go"
	refBinary2062F = "ref2062F.bin"
	maxTests       = 160
	maxTotalN2     = 3500000
)

type testCase struct {
	n int
	a []int64
	b []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch at output %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2062F, refSource2062F)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2062F), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, tests []testCase) ([]int64, error) {
	total := 0
	for _, tc := range tests {
		total += tc.n - 1
	}
	fields := strings.Fields(out)
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d numbers, got %d", total, len(fields))
	}
	res := make([]int64, total)
	for i, tok := range fields {
		var v int64
		if _, err := fmt.Sscan(tok, &v); err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i := 0; i < tc.n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", tc.a[i], tc.b[i])
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2062))
	var tests []testCase
	totalN2 := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN2 += tc.n * tc.n
	}

	add(testCase{
		n: 2,
		a: []int64{0, 0},
		b: []int64{0, 0},
	})
	add(testCase{
		n: 3,
		a: []int64{0, 2, 3},
		b: []int64{2, 1, 5},
	})
	add(testCase{
		n: 4,
		a: []int64{1, 7, 6, 1},
		b: []int64{8, 9, 9, 1},
	})
	add(testCase{
		n: 5,
		a: []int64{1000000000, 0, 500000000, 123456789, 987654321},
		b: []int64{999999999, 1000000000, 400000000, 123456789, 5},
	})

	for len(tests) < maxTests {
		remain := maxTotalN2 - totalN2
		if remain <= 0 {
			break
		}

		var n int
		switch rnd.Intn(5) {
		case 0:
			n = rnd.Intn(5) + 2
		case 1:
			n = rnd.Intn(15) + 5
		case 2:
			n = rnd.Intn(60) + 20
		case 3:
			n = rnd.Intn(250) + 80
		default:
			n = rnd.Intn(700) + 120
		}
		if n < 2 {
			n = 2
		}
		if n*n > remain {
			if remain < 4 {
				break
			}
			n = int(randRange(int64(2), int64(minInt(remain/2, 1500)), rnd))
			if n < 2 {
				n = 2
			}
		}

		aVals := make([]int64, n)
		bVals := make([]int64, n)
		for i := 0; i < n; i++ {
			mode := rnd.Intn(6)
			switch mode {
			case 0:
				val := randRange(0, 10, rnd)
				aVals[i] = val
				bVals[i] = val
			case 1:
				base := randRange(0, 1000000000, rnd)
				diff := randRange(0, 200, rnd)
				aVals[i] = clamp1e9(base)
				bVals[i] = clamp1e9(base + diff)
			case 2:
				base := randRange(0, 1000000000, rnd)
				diff := randRange(0, 200, rnd)
				aVals[i] = clamp1e9(base + diff)
				bVals[i] = clamp1e9(base)
			default:
				aVals[i] = randRange(0, 1000000000, rnd)
				bVals[i] = randRange(0, 1000000000, rnd)
			}
		}

		add(testCase{
			n: n,
			a: aVals,
			b: bVals,
		})
	}

	return tests
}

func randRange(lo, hi int64, rnd *rand.Rand) int64 {
	if hi < lo {
		lo, hi = hi, lo
	}
	return lo + rnd.Int63n(hi-lo+1)
}

func clamp1e9(x int64) int64 {
	if x < 0 {
		return 0
	}
	if x > 1000000000 {
		return 1000000000
	}
	return x
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
