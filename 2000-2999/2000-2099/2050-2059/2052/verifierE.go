package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource = "2052E.go"
	refBinary = "ref2052E.bin"
)

type testCase struct {
	name  string
	eqStr string
}

func main() {
	if len(os.Args) != 2 {
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

	tests := buildTests()
	for idx, tc := range tests {
		input := []byte(tc.eqStr + "\n")

		// Sanity check reference output passes validator.
		if refOut, err := runProgram(refPath, input); err != nil || validateAnswer(tc.eqStr, refOut) != nil {
			fmt.Printf("reference failed on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			if err == nil {
				fmt.Println("Reference output:")
				fmt.Println(refOut)
			}
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}

		if err := validateAnswer(tc.eqStr, candOut); err != nil {
			fmt.Printf("test %d (%s) failed: %v\n", idx+1, tc.name, err)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample-correct", eqStr: "2+2=4"},
		{name: "sample-fixable", eqStr: "123456789+9876543210=111111110+11-1"},
		{name: "sample-impossible", eqStr: "10+9=10"},
		{name: "move-leading", eqStr: "24=55-13"},
		{name: "large-impossible", eqStr: "1000000000-10=9999999999"},
	}

	// Create fixable cases by shuffling a digit in a correct equality.
	fixables := []string{
		"42=21+21",
		"33333=11111+22222",
		"1111111111=555555555+555555556",
	}
	for i, eq := range fixables {
		tests = append(tests, testCase{
			name:  fmt.Sprintf("fixable-%d", i+1),
			eqStr: scrambleOneDigit(eq),
		})
	}

	// Already correct but more complex.
	tests = append(tests, testCase{name: "multi-term-correct", eqStr: "100+200-50=250"})

	// Random tests; all numbers valid (length <=10, no leading zeros).
	rnd := rand.New(rand.NewSource(20240530))
	for i := 0; i < 40; i++ {
		eq := randomEquality(rnd)
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random-%d", i+1),
			eqStr: eq,
		})
	}

	return tests
}

func scrambleOneDigit(eq string) string {
	// assumes eq is correct; move one digit to make likely incorrect but fixable.
	positions := make([]int, 0, len(eq))
	for i := range eq {
		if eq[i] >= '0' && eq[i] <= '9' {
			positions = append(positions, i)
		}
	}
	if len(positions) == 0 {
		return eq
	}
	fromIdx := positions[rand.Intn(len(positions))]
	to := rand.Intn(len(eq))
	if to >= fromIdx {
		to++
	}
	return moveDigit(eq, fromIdx, to)
}

func randomNumber(rnd *rand.Rand) string {
	length := rnd.Intn(10) + 1
	var sb strings.Builder
	first := rnd.Intn(9) + 1
	sb.WriteByte(byte('0' + first))
	for i := 1; i < length; i++ {
		sb.WriteByte(byte('0' + rnd.Intn(10)))
	}
	return sb.String()
}

func randomExpression(rnd *rand.Rand, maxLen int) string {
	parts := make([]string, 0)
	curLen := 0
	for len(parts) == 0 || curLen+2 < maxLen {
		num := randomNumber(rnd)
		if curLen+len(num) > maxLen {
			if len(parts) == 0 {
				return num
			}
			break
		}
		parts = append(parts, num)
		curLen += len(num)
		if curLen+1 >= maxLen {
			break
		}
		if rnd.Intn(2) == 0 {
			parts = append(parts, "+")
		} else {
			parts = append(parts, "-")
		}
		curLen++
	}
	if len(parts) > 0 {
		if parts[len(parts)-1] == "+" || parts[len(parts)-1] == "-" {
			parts = parts[:len(parts)-1]
		}
	}
	return strings.Join(parts, "")
}

func randomEquality(rnd *rand.Rand) string {
	left := randomExpression(rnd, 45)
	right := randomExpression(rnd, 45)
	return left + "=" + right
}

func parseExpr(s string) (bool, *big.Int) {
	if len(s) == 0 {
		return false, nil
	}
	sum := big.NewInt(0)
	sign := int64(1)
	i := 0
	for {
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if start == i {
			return false, nil
		}
		if s[start] == '0' && i-start > 1 {
			return false, nil
		}
		if i-start > 10 {
			return false, nil
		}
		var num int64
		for k := start; k < i; k++ {
			num = num*10 + int64(s[k]-'0')
		}
		term := big.NewInt(num)
		if sign == 1 {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}

		if i == len(s) {
			break
		}
		if s[i] != '+' && s[i] != '-' {
			return false, nil
		}
		if s[i] == '+' {
			sign = 1
		} else {
			sign = -1
		}
		i++
		if i == len(s) {
			return false, nil
		}
	}
	return true, sum
}

func parseEquality(s string) (bool, *big.Int, *big.Int) {
	eqPos := strings.IndexByte(s, '=')
	if eqPos <= 0 || eqPos == len(s)-1 {
		return false, nil, nil
	}
	if strings.Count(s, "=") != 1 {
		return false, nil, nil
	}
	lok, lval := parseExpr(s[:eqPos])
	rok, rval := parseExpr(s[eqPos+1:])
	if !lok || !rok {
		return false, nil, nil
	}
	return true, lval, rval
}

func isCorrectEquality(s string) bool {
	ok, l, r := parseEquality(s)
	return ok && l.Cmp(r) == 0
}

func canFixWithOneMove(orig string) bool {
	for i := 0; i < len(orig); i++ {
		if orig[i] < '0' || orig[i] > '9' {
			continue
		}
		for to := 0; to <= len(orig); to++ {
			if to == i {
				continue
			}
			if res := moveDigit(orig, i, to); isCorrectEquality(res) {
				return true
			}
		}
	}
	return false
}

func moveDigit(s string, from, to int) string {
	ch := s[from]
	t := s[:from] + s[from+1:]
	if to < 0 {
		to = 0
	}
	if to > len(t) {
		to = len(t)
	}
	return t[:to] + string(ch) + t[to:]
}

func obtainedBySingleDigitMove(orig, cand string) bool {
	if len(orig) != len(cand) {
		return false
	}
	for i := 0; i < len(orig); i++ {
		if orig[i] < '0' || orig[i] > '9' {
			continue
		}
		for to := 0; to <= len(orig); to++ {
			if to == i {
				continue
			}
			if moveDigit(orig, i, to) == cand {
				return true
			}
		}
	}
	return false
}

func validateAnswer(inputEq, output string) error {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return fmt.Errorf("expected single token, got %d", len(fields))
	}
	out := fields[0]

	origCorrect := isCorrectEquality(inputEq)
	fixable := canFixWithOneMove(inputEq)

	if out == "Correct" {
		if !origCorrect {
			return fmt.Errorf("output Correct but equality is not correct")
		}
		return nil
	}
	if out == "Impossible" {
		if origCorrect || fixable {
			return fmt.Errorf("output Impossible but correction exists")
		}
		return nil
	}

	if !isCorrectEquality(out) {
		return fmt.Errorf("provided equality is not correct: %s", out)
	}
	if !obtainedBySingleDigitMove(inputEq, out) {
		return fmt.Errorf("output is not obtainable by moving one digit")
	}
	return nil
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Print(string(in))
}
