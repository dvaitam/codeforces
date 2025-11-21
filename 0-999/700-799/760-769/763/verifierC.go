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

const refSource = "0-999/700-799/760-769/763/763C.go"

type testCase struct {
	name  string
	input string
	n     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		expLine := strings.TrimSpace(refOut)
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := verifyOutput(tc, strings.TrimSpace(candOut), expLine); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, expLine, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-763C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref763C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func verifyOutput(tc testCase, candLine string, refLine string) error {
	if refLine == "-1" {
		if candLine != "-1" {
			return fmt.Errorf("expected -1, got %s", candLine)
		}
		return nil
	}
	var x, d int64
	fields := strings.Fields(candLine)
	if len(fields) != 2 {
		return fmt.Errorf("expected two integers, got %q", candLine)
	}
	mod, n, seq := parseInput(tc.input)
	var err error
	x, err = strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid x: %v", err)
	}
	d, err = strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid d: %v", err)
	}
	if x < 0 || x >= mod || d < 0 || d >= mod {
		return fmt.Errorf("x or d out of range")
	}
	if !checkSequence(seq, mod, n, x, d) {
		return fmt.Errorf("candidate progression invalid")
	}
	return nil
}

func parseInput(input string) (mod int64, n int64, seq []int64) {
	lines := strings.Fields(input)
	modVal, _ := strconv.ParseInt(lines[0], 10, 64)
	nVal, _ := strconv.ParseInt(lines[1], 10, 64)
	arr := make([]int64, nVal)
	for i := int64(0); i < nVal; i++ {
		v, _ := strconv.ParseInt(lines[2+i], 10, 64)
		arr[i] = v % modVal
	}
	return modVal, nVal, arr
}

func checkSequence(seq []int64, mod int64, n int64, x int64, d int64) bool {
	if n == 0 {
		return true
	}
	if n == 1 {
		return seq[0] == x
	}
	seen := make(map[int64]bool, n)
	for i := int64(0); i < n; i++ {
		val := (x + d*i) % mod
		seen[val] = true
	}
	if int64(len(seen)) != n {
		return false
	}
	for _, v := range seq {
		if !seen[v] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		formatCase("single", 7, []int{3}),
		formatCase("arithmetic_progression", 17, []int{0, 5, 10, 15}),
		formatCase("full_range", 7, []int{0, 1, 2, 3, 4, 5, 6}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatCase(name string, mod int64, seq []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", mod, len(seq))
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String(), n: len(seq)}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	mod := bigPrime(rng)
	n := rng.Intn(100000) + 1
	seq := randomDistinct(rng, mod, n)
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: formatInput(mod, seq),
		n:     n,
	}
}

func stressCase() testCase {
	mod := int64(1_000_000_007)
	n := 100000
	seq := make([]int, n)
	cur := rand.New(rand.NewSource(42))
	used := make(map[int]int)
	for i := 0; i < n; i++ {
		for {
			val := cur.Intn(int(mod))
			if _, ok := used[val]; !ok {
				seq[i] = val
				used[val] = 1
				break
			}
		}
	}
	return testCase{
		name:  "stress_max",
		input: formatInput(mod, seq),
		n:     n,
	}
}

func formatInput(mod int64, seq []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", mod, len(seq))
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomDistinct(rng *rand.Rand, mod int64, n int) []int {
	if mod <= int64(n) {
		mod = int64(n) + 1
	}
	set := make(map[int]int, n)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			val := rng.Intn(int(mod))
			if _, ok := set[val]; !ok {
				set[val] = 1
				res[i] = val
				break
			}
		}
	}
	return res
}

func bigPrime(rng *rand.Rand) int64 {
	primes := []int64{
		1000003, 1000033, 1000037, 1000039,
		1000000007, 1000000009, 1000000033,
	}
	return primes[rng.Intn(len(primes))]
}
