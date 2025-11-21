package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n       int
	q       int
	parent  []int
	p       []int
	queries [][2]int
}

const (
	maxTotalN = 300000
	maxTotalQ = 100000
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)
	totalQueries := totalQ(tests)

	expected, err := runAndParse(refBin, input, totalQueries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalQueries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for tIdx, tc := range tests {
		for qIdx := 0; qIdx < tc.q; qIdx++ {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d query %d: expected %s got %s\n", tIdx+1, qIdx+1, expected[idx], got[idx])
				fmt.Fprintf(os.Stderr, "n=%d q=%d\nparent=%v\np=%v\nquery=%v\n", tc.n, tc.q, tc.parent, tc.p, tc.queries[qIdx])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2002D2.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2002D2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, count int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != count {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", count, len(fields), out)
	}
	res := make([]string, count)
	for i, f := range fields {
		res[i] = strings.ToUpper(f)
	}
	return res, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	sumN := totalN(tests)
	sumQ := totalQ(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for sumN < maxTotalN && sumQ < maxTotalQ {
		n := rng.Intn(min(2000, maxTotalN-sumN)) + 2
		q := rng.Intn(min(2000, maxTotalQ-sumQ)) + 2
		parent := make([]int, n+1)
		for i := 2; i <= n; i++ {
			parent[i] = rng.Intn(i-1) + 1
		}
		p := randPermutation(rng, n)
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(n-1) + 1
			if y >= x {
				y++
			}
			queries[i] = [2]int{x, y}
		}
		tests = append(tests, testCase{
			n:       n,
			q:       q,
			parent:  parent,
			p:       p,
			queries: queries,
		})
		sumN += n
		sumQ += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:       3,
			q:       3,
			parent:  []int{0, 0, 1, 1},
			p:       []int{0, 1, 2, 3},
			queries: [][2]int{{1, 2}, {2, 3}, {1, 3}},
		},
	}
}

func randPermutation(rng *rand.Rand, n int) []int {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	rng.Shuffle(n, func(i, j int) {
		p[i+1], p[j+1] = p[j+1], p[i+1]
	})
	return p
}

func totalN(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	return total
}

func totalQ(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += tc.q
	}
	return total
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i := 2; i <= tc.n; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.parent[i]))
		}
		sb.WriteByte('\n')
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.p[i]))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
		}
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
