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
	maxValue = 100000000
	maxX     = int64(1_000_000_000_000_000_000)
)

type caseSpec struct {
	n, k int
	x    int64
	arr  []int
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2080-2089", "2086")
	tmp, err := os.CreateTemp("", "ref2086B")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2086B.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func normalizeOutput(s string) string {
	return strings.TrimSpace(s)
}

func buildInput(specs []caseSpec) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(specs)))
	for _, cs := range specs {
		x := clampX(cs.x)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", cs.n, cs.k, x))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func uniformCase(n, k, val int, x int64) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, clampX(x)))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func clampX(x int64) int64 {
	if x < 1 {
		return 1
	}
	if x > maxX {
		return maxX
	}
	return x
}

func randRange64(rng *rand.Rand, lo, hi int64) int64 {
	if hi < lo {
		hi = lo
	}
	return rng.Int63n(hi-lo+1) + lo
}

func pickSize(rng *rand.Rand, maxVal int) int {
	if maxVal <= 1 {
		return 1
	}
	candidates := []int{
		min(maxVal, 5),
		min(maxVal, 50),
		min(maxVal, 500),
		min(maxVal, 5000),
		maxVal,
	}
	choice := candidates[rng.Intn(len(candidates))]
	if choice < 1 {
		choice = 1
	}
	return rng.Intn(choice) + 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	nBudget := 200000
	kBudget := 200000
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		remain := t - i - 1
		maxNAllowed := max(1, nBudget-remain)
		maxNAllowed = min(100000, maxNAllowed)
		n := pickSize(rng, maxNAllowed)
		nBudget -= n

		maxKAllowed := max(1, kBudget-remain)
		maxKAllowed = min(100000, maxKAllowed)
		k := pickSize(rng, maxKAllowed)
		kBudget -= k

		arr := make([]int, n)
		var total int64
		for j := 0; j < n; j++ {
			val := rng.Intn(maxValue) + 1
			arr[j] = val
			total += int64(val)
		}
		maxSum := total * int64(k)
		if maxSum < 1 {
			maxSum = 1
		}
		var x int64
		switch rng.Intn(5) {
		case 0:
			x = 1
		case 1:
			x = maxSum
		case 2:
			x = maxSum + int64(rng.Intn(1000)+1)
		case 3:
			x = randRange64(rng, 1, maxSum)
		default:
			x = randRange64(rng, 1, maxX)
		}
		x = clampX(x)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, x))
		for j, val := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func fixedTests() []string {
	tests := []string{
		buildInput([]caseSpec{
			{n: 1, k: 1, x: 1, arr: []int{1}},
			{n: 2, k: 3, x: 5, arr: []int{1, 2}},
		}),
		buildInput([]caseSpec{
			{n: 5, k: 3, x: 10, arr: []int{3, 4, 2, 1, 5}},
			{n: 4, k: 2, x: 9, arr: []int{1, 2, 3, 4}},
		}),
	}
	tests = append(tests,
		uniformCase(1, 100000, 1, 1),
		uniformCase(1, 100000, 1, 200000),
		uniformCase(100000, 1, 1, 50000),
		uniformCase(100000, 100000, maxValue, maxX),
	)
	return tests
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := fixedTests()
	for len(tests) < 120 {
		tests = append(tests, generateRandomInput(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, input := range tests {
		expect, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Fprintf(os.Stderr, "input:\n%s", input)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n",
				idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
