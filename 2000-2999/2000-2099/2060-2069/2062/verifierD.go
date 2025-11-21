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
	refSource2062D = "2062D.go"
	refBinary2062D = "ref2062D.bin"
	maxTests       = 120
	maxTotalN      = 200000
)

type testCase struct {
	n     int
	l     []int64
	r     []int64
	edges [][2]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2062D, refSource2062D)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2062D), nil
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
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i := 0; i < tc.n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", tc.l[i], tc.r[i])
		}
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2062))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	add(testCase{
		n:     1,
		l:     []int64{0},
		r:     []int64{10},
		edges: nil,
	})
	add(testCase{
		n:     2,
		l:     []int64{0, 5},
		r:     []int64{5, 5},
		edges: [][2]int{{1, 2}},
	})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxN := 5000
		if remain < maxN {
			maxN = remain
		}
		n := rnd.Intn(maxN-1) + 1
		if n == 0 {
			n = 1
		}

		lVals := make([]int64, n)
		rVals := make([]int64, n)
		for i := 0; i < n; i++ {
			if rnd.Intn(5) == 0 {
				val := int64(rnd.Intn(20))
				lVals[i] = val
				rVals[i] = val
			} else {
				a := int64(rnd.Intn(1000))
				b := a + int64(rnd.Intn(1000))
				lVals[i] = a
				rVals[i] = b
			}
			if rnd.Intn(10) == 0 {
				lVals[i] = int64(rnd.Int31())
				rVals[i] = lVals[i] + int64(rnd.Intn(1000000000))
			}
		}

		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			u := rnd.Intn(v-1) + 1
			edges = append(edges, [2]int{u, v})
		}

		add(testCase{
			n:     n,
			l:     lVals,
			r:     rVals,
			edges: edges,
		})
	}
	return tests
}
