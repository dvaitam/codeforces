package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 1000000007

type testCase struct {
	input string
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswer(refOut)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		if ans != exp {
			fmt.Printf("Wrong answer on test %d (%s): expected %d got %d\nInput:\n%s", i+1, tc.desc, exp, ans, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1045H.bin"
	cmd := exec.Command("go", "build", "-o", path, "1045H.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswer(out string) (int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var ans int64
	if _, err := fmt.Fscan(reader, &ans); err != nil {
		return 0, fmt.Errorf("failed to read answer: %v", err)
	}
	if ans < 0 || ans >= mod {
		return 0, fmt.Errorf("answer %d outside [0,%d)", ans, mod)
	}
	rest := strings.TrimSpace(readRemaining(reader))
	if rest != "" {
		return 0, fmt.Errorf("unexpected extra output: %q", rest)
	}
	return int(ans), nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		s, err := r.ReadString('\n')
		sb.WriteString(s)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc, input string) {
		tests = append(tests, testCase{desc: desc, input: input})
	}

	// Statement sample
	add("statement-sample", formatInput("10", "1001", 0, 0, 1, 1))
	add("single-number", formatInput("1", "1", 0, 0, 0, 0))

	target := "110"
	c00, c01, c10, c11 := substringCounts(target)
	add("tight-interval", formatInput("110", "110", c00, c01, c10, c11))

	// Random positive cases built from actual strings
	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 35 {
		s := randomBinaryString(rng, 1, 25)
		if s == "" {
			continue
		}
		aStr, bStr := intervalAroundString(rng, s)
		c00, c01, c10, c11 := substringCounts(s)
		desc := fmt.Sprintf("random-positive-%d", len(tests))
		add(desc, formatInput(aStr, bStr, c00, c01, c10, c11))
	}

	// Random negative / arbitrary counts
	for len(tests) < 60 {
		aStr, bStr := randomInterval(rng, 1, 30)
		nc00 := rng.Intn(100001)
		nc01 := rng.Intn(100001)
		nc10 := rng.Intn(100001)
		nc11 := rng.Intn(100001)
		desc := fmt.Sprintf("random-negative-%d", len(tests))
		add(desc, formatInput(aStr, bStr, nc00, nc01, nc10, nc11))
	}

	// Large structured case
	large := buildLargeString(800, 400)
	lc00, lc01, lc10, lc11 := substringCounts(large)
	add("large-equal", formatInput(large, large, lc00, lc01, lc10, lc11))

	// Wide interval with strict counts (infeasible)
	bigB := "1" + strings.Repeat("0", 2000)
	add("wide-zero", formatInput("1", bigB, 100000, 100000, 100000, 100000))

	return tests
}

func formatInput(a, b string, c00, c01, c10, c11 int) string {
	return fmt.Sprintf("%s\n%s\n%d\n%d\n%d\n%d\n", a, b, c00, c01, c10, c11)
}

func substringCounts(s string) (int, int, int, int) {
	var c00, c01, c10, c11 int
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '0' && s[i+1] == '0' {
			c00++
		} else if s[i] == '0' && s[i+1] == '1' {
			c01++
		} else if s[i] == '1' && s[i+1] == '0' {
			c10++
		} else if s[i] == '1' && s[i+1] == '1' {
			c11++
		}
	}
	return c00, c01, c10, c11
}

func randomBinaryString(rng *rand.Rand, minLen, maxLen int) string {
	if maxLen < minLen {
		maxLen = minLen
	}
	length := rng.Intn(maxLen-minLen+1) + minLen
	var sb strings.Builder
	sb.Grow(length)
	sb.WriteByte('1')
	for i := 1; i < length; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func intervalAroundString(rng *rand.Rand, s string) (string, string) {
	val := mustParseBinary(s)
	left := rng.Intn(6)
	right := rng.Intn(6)
	aVal := new(big.Int).Sub(val, big.NewInt(int64(left)))
	if aVal.Sign() <= 0 {
		aVal.SetInt64(1)
	}
	bVal := new(big.Int).Add(val, big.NewInt(int64(right)))
	if bVal.Cmp(aVal) < 0 {
		bVal.Set(aVal)
	}
	return aVal.Text(2), bVal.Text(2)
}

func randomInterval(rng *rand.Rand, minLen, maxLen int) (string, string) {
	bStr := randomBinaryString(rng, minLen, maxLen)
	bVal := mustParseBinary(bStr)
	if bVal.Sign() == 0 {
		bVal.SetInt64(1)
	}
	r := new(big.Int).Rand(rng, bVal)
	r.Add(r, big.NewInt(1))
	return r.Text(2), bVal.Text(2)
}

func mustParseBinary(s string) *big.Int {
	val := new(big.Int)
	if _, ok := val.SetString(s, 2); !ok {
		panic("invalid binary string")
	}
	return val
}

func buildLargeString(zeros, ones int) string {
	var sb strings.Builder
	sb.Grow(1 + zeros + ones)
	sb.WriteByte('1')
	for i := 0; i < zeros; i++ {
		sb.WriteByte('0')
	}
	for i := 0; i < ones; i++ {
		sb.WriteByte('1')
	}
	return sb.String()
}
