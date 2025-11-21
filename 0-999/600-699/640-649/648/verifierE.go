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

type testCase struct {
	input    []byte
	k        int
	valueSet map[int]struct{}
	maxLen   int
}

type solution struct {
	possible bool
	number   string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "648E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseSolution(out string) (solution, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return solution{}, fmt.Errorf("empty output")
	}
	status := strings.ToUpper(fields[0])
	switch status {
	case "NO":
		return solution{possible: false}, nil
	case "YES":
		if len(fields) < 2 {
			return solution{}, fmt.Errorf("YES without number")
		}
		return solution{possible: true, number: fields[1]}, nil
	default:
		return solution{}, fmt.Errorf("unexpected verdict %q", fields[0])
	}
}

func digitLen(x int) int {
	if x == 0 {
		return 1
	}
	cnt := 0
	for v := x; v > 0; v /= 10 {
		cnt++
	}
	return cnt
}

func newTestCase(values []int, k int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(values), k)
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	valueSet := make(map[int]struct{}, len(values))
	maxLen := 1
	for _, v := range values {
		valueSet[v] = struct{}{}
		if l := digitLen(v); l > maxLen {
			maxLen = l
		}
	}
	return testCase{
		input:    []byte(sb.String()),
		k:        k,
		valueSet: valueSet,
		maxLen:   maxLen,
	}
}

func deterministicTests() []testCase {
	tests := []testCase{
		newTestCase([]int{1}, 1),
		newTestCase([]int{2}, 3),
		newTestCase([]int{0}, 5),
		newTestCase([]int{5, 55, 555}, 5),
		newTestCase([]int{12, 34, 56}, 7),
		newTestCase([]int{101, 10, 1}, 3),
	}
	big := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		big[i] = i
	}
	tests = append(tests, newTestCase(big, 997))
	return tests
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%10 == 0:
			n = rnd.Intn(1000) + 1
		case i%4 == 0:
			n = rnd.Intn(200) + 1
		default:
			n = rnd.Intn(50) + 1
		}
		k := rnd.Intn(1000) + 1
		values := make([]int, n)
		for j := 0; j < n; j++ {
			values[j] = rnd.Intn(1_000_000_000 + 1)
		}
		if n > 0 && rnd.Intn(5) == 0 {
			values[rnd.Intn(n)] = 0
		}
		tests = append(tests, newTestCase(values, k))
	}
	return tests
}

func checkNumber(tc testCase, num string) error {
	if !tc.valueSetValid() {
		return fmt.Errorf("internal error: empty value set")
	}
	if len(num) == 0 {
		return fmt.Errorf("empty number")
	}
	if num[0] == '0' && len(num) > 1 {
		return fmt.Errorf("leading zeros in %q", num)
	}
	if tc.k <= 0 {
		return fmt.Errorf("invalid k %d", tc.k)
	}
	mod := 0
	for _, ch := range num {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("non-digit character %q", ch)
		}
		mod = (mod*10 + int(ch-'0')) % tc.k
	}
	if mod%tc.k != 0 {
		return fmt.Errorf("number %q not divisible by %d", num, tc.k)
	}
	if !canForm(num, tc) {
		return fmt.Errorf("number %q cannot be formed from given pieces", num)
	}
	return nil
}

func (tc testCase) valueSetValid() bool {
	return len(tc.valueSet) > 0
}

func canForm(num string, tc testCase) bool {
	dp := make([]bool, len(num)+1)
	dp[0] = true
	for i := 0; i < len(num); i++ {
		if !dp[i] {
			continue
		}
		if num[i] == '0' {
			if _, ok := tc.valueSet[0]; ok {
				dp[i+1] = true
			}
			continue
		}
		val := 0
		for l := 1; l <= tc.maxLen && i+l <= len(num); l++ {
			val = val*10 + int(num[i+l-1]-'0')
			if val > 1_000_000_000 {
				break
			}
			if _, ok := tc.valueSet[val]; ok {
				dp[i+l] = true
			}
		}
	}
	return dp[len(num)]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(120)...)

	for idx, tc := range tests {
		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expSol, err := parseSolution(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotSol, err := parseSolution(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expSol.possible != gotSol.possible {
			if expSol.possible {
				fmt.Fprintf(os.Stderr, "case %d: expected YES but got NO\n", idx+1)
			} else {
				fmt.Fprintf(os.Stderr, "case %d: expected NO but got YES\n", idx+1)
			}
			os.Exit(1)
		}
		if !expSol.possible {
			continue
		}
		if len(gotSol.number) != len(expSol.number) {
			fmt.Fprintf(os.Stderr, "case %d: expected length %d got %d\n", idx+1, len(expSol.number), len(gotSol.number))
			os.Exit(1)
		}
		if err := checkNumber(tc, gotSol.number); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
