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
)

const (
	refSource2133B = "2133B.go"
	refBinary2133B = "ref2133B.bin"
	maxTests       = 160
	maxTotalN      = 180000
)

type testCase struct {
	n int
	g []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2133B, refSource2133B)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2133B), nil
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

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
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
		for i, v := range tc.g {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2133))
	var tests []testCase
	totalN := 0

	add := func(arr []int64) {
		if len(arr) < 2 {
			return
		}
		if totalN+len(arr) > maxTotalN || len(tests) >= maxTests {
			return
		}
		tests = append(tests, testCase{n: len(arr), g: append([]int64(nil), arr...)})
		totalN += len(arr)
	}

	// Deterministic small cases
	add([]int64{1, 2})
	add([]int64{2, 1, 5, 2})
	add([]int64{1_000_000_000, 1_000_000_000, 1_000_000_000, 1_000_000_000})
	add([]int64{6, 3, 1, 4, 1, 5, 9})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		n := rnd.Intn(1200) + 2
		if n > remain {
			n = remain
		}
		if n < 2 {
			break
		}
		g := make([]int64, n)
		mode := rnd.Intn(4)
		for i := 0; i < n; i++ {
			switch mode {
			case 0:
				g[i] = int64(rnd.Intn(5) + 1)
			case 1:
				g[i] = int64(rnd.Intn(1_000) + 1)
			case 2:
				g[i] = int64(rnd.Intn(1_000_000) + 1)
			default:
				g[i] = rnd.Int63n(1_000_000_000) + 1
			}
		}
		add(g)
	}
	return tests
}
