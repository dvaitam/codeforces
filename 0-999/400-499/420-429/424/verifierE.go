package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceE = "424E.go"
	refBinaryE = "ref424E.bin"
	totalTests = 120
	tolerance  = 1e-6
)

type testCase struct {
	n      int
	levels []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if !closeEnough(refVal, candVal) {
			fmt.Printf("test %d failed: expected %.10f, got %.10f\n", idx+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryE, refSourceE)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryE), nil
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

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, level := range tc.levels {
		sb.WriteString(level)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func parseOutput(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite value %v", val)
	}
	return val, nil
}

func closeEnough(expected, actual float64) bool {
	diff := math.Abs(expected - actual)
	allowed := tolerance * math.Max(1.0, math.Abs(expected))
	return diff <= allowed+1e-12
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 2, levels: []string{"RGB", "RGB"}},
		{n: 3, levels: []string{"RRR", "GGG", "BBB"}},
		{n: 4, levels: []string{"RRG", "GBB", "BRG", "GRB"}},
	}

	colors := []byte{'R', 'G', 'B'}
	// exhaustive for n=2 over a subset
	for mask := 0; mask < 27 && len(tests) < 30; mask++ {
		cur := mask
		level := make([]byte, 3)
		for i := 0; i < 3; i++ {
			level[i] = colors[cur%3]
			cur /= 3
		}
		tests = append(tests, testCase{n: 2, levels: []string{string(level), reverseLevel(level)}})
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(5) + 2 // 2..6
		levels := make([]string, n)
		for i := 0; i < n; i++ {
			levels[i] = randomLevel(rnd)
		}
		tests = append(tests, testCase{n: n, levels: levels})
	}

	tests = append(tests, buildLayeredTower([]string{"RGB", "RGB", "RGB", "RGB", "RGB", "RGB"}))
	tests = append(tests, buildLayeredTower([]string{"RRR", "RRR", "RRR", "RRR", "RRR", "RRR"}))
	tests = append(tests, buildLayeredTower([]string{"RRG", "RRG", "RRG", "BBG", "BBG", "BBG"}))
	tests = append(tests, buildLayeredTower([]string{"RGB", "BGR", "GBR", "RBG", "BRG", "GRB"}))
	tests = append(tests, buildLayeredTower([]string{"BBB", "GGG", "RRR", "BBB", "GGG", "RRR"}))

	return tests
}

func reverseLevel(level []byte) string {
	tmp := make([]byte, len(level))
	copy(tmp, level)
	for i := 0; i < len(tmp)/2; i++ {
		tmp[i], tmp[len(tmp)-1-i] = tmp[len(tmp)-1-i], tmp[i]
	}
	return string(tmp)
}

func randomLevel(rnd *rand.Rand) string {
	colors := []byte{'R', 'G', 'B'}
	b := make([]byte, 3)
	for i := 0; i < 3; i++ {
		b[i] = colors[rnd.Intn(3)]
	}
	return string(b)
}

func buildLayeredTower(levels []string) testCase {
	n := len(levels)
	if n > 6 {
		n = 6
		levelCopy := make([]string, n)
		copy(levelCopy, levels[:n])
		return testCase{n: n, levels: levelCopy}
	}
	levelCopy := make([]string, n)
	copy(levelCopy, levels)
	return testCase{n: n, levels: levelCopy}
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
