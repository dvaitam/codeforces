package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	input string
	p     int64
	q     int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateExpression(strings.TrimSpace(refOut), tc.p, tc.q); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateExpression(strings.TrimSpace(gotOut), tc.p, tc.q); err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1912E.go",
		filepath.Join("1000-1999", "1900-1999", "1910-1919", "1912", "1912E.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1912E.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1912E_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generateTests() []testCase {
	tests := []testCase{
		{input: "1998 -3192\n", p: 1998, q: -3192}, // sample
		{input: "413 908\n", p: 413, q: 908},       // sample
		{input: "0 0\n", p: 0, q: 0},
		{input: "1 -1\n", p: 1, q: -1},
		{input: "-1 1\n", p: -1, q: 1},
		{input: "1234567890 -1234567890\n", p: 1234567890, q: -1234567890},
		{input: "-1000000000000000000 1000000000000000000\n", p: -1_000_000_000_000_000_000, q: 1_000_000_000_000_000_000},
		{input: "1000000000000000000 -1000000000000000000\n", p: 1_000_000_000_000_000_000, q: -1_000_000_000_000_000_000},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		p := randRange(rng, -1_000_000_000_000_000_000, 1_000_000_000_000_000_000)
		q := randRange(rng, -1_000_000_000_000_000_000, 1_000_000_000_000_000_000)
		tests = append(tests, testCase{
			input: fmt.Sprintf("%d %d\n", p, q),
			p:     p,
			q:     q,
		})
	}
	return tests
}

func randRange(r *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + r.Int63n(hi-lo+1)
}

func validateExpression(expr string, p, q int64) error {
	if expr == "" {
		return fmt.Errorf("empty output")
	}
	if len(expr) > 1000 {
		return fmt.Errorf("expression too long: %d", len(expr))
	}
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if (c < '0' || c > '9') && c != '+' && c != '-' && c != '*' {
			return fmt.Errorf("invalid character %q", c)
		}
	}
	pBig := big.NewInt(p)
	qBig := big.NewInt(q)

	forward, err := parseAndEval(expr)
	if err != nil {
		return fmt.Errorf("invalid forward expression: %v", err)
	}
	if forward.Cmp(pBig) != 0 {
		return fmt.Errorf("forward value mismatch: got %s expected %s", forward.String(), pBig.String())
	}

	reversed := reverseString(expr)
	backward, err := parseAndEval(reversed)
	if err != nil {
		return fmt.Errorf("invalid reversed expression: %v", err)
	}
	if backward.Cmp(qBig) != 0 {
		return fmt.Errorf("reversed value mismatch: got %s expected %s", backward.String(), qBig.String())
	}
	return nil
}

func parseAndEval(expr string) (*big.Int, error) {
	numbers, ops, err := tokenize(expr)
	if err != nil {
		return nil, err
	}

	vals := make([]*big.Int, 0, len(numbers))
	opStack := make([]byte, 0, len(ops))

	for i, numStr := range numbers {
		val, ok := new(big.Int).SetString(numStr, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse number %q", numStr)
		}
		vals = append(vals, val)
		if i < len(ops) {
			curOp := ops[i]
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(curOp) {
				if err := applyOp(&vals, &opStack); err != nil {
					return nil, err
				}
			}
			opStack = append(opStack, curOp)
		}
	}

	for len(opStack) > 0 {
		if err := applyOp(&vals, &opStack); err != nil {
			return nil, err
		}
	}
	if len(vals) != 1 {
		return nil, fmt.Errorf("evaluation error")
	}
	return vals[0], nil
}

func tokenize(expr string) ([]string, []byte, error) {
	if expr == "" {
		return nil, nil, fmt.Errorf("empty expression")
	}
	var numbers []string
	var ops []byte
	for i := 0; i < len(expr); {
		if expr[i] < '0' || expr[i] > '9' {
			return nil, nil, fmt.Errorf("unexpected character %q", expr[i])
		}
		start := i
		for i < len(expr) && expr[i] >= '0' && expr[i] <= '9' {
			i++
		}
		num := expr[start:i]
		if len(num) > 1 && num[0] == '0' {
			return nil, nil, fmt.Errorf("leading zero in %q", num)
		}
		numbers = append(numbers, num)
		if i == len(expr) {
			break
		}
		op := expr[i]
		if op != '+' && op != '-' && op != '*' {
			return nil, nil, fmt.Errorf("invalid operator %q", op)
		}
		ops = append(ops, op)
		i++
		if i == len(expr) {
			return nil, nil, fmt.Errorf("expression ends with operator")
		}
	}
	if len(numbers) != len(ops)+1 {
		return nil, nil, fmt.Errorf("malformed expression")
	}
	return numbers, ops, nil
}

func precedence(op byte) int {
	if op == '*' {
		return 2
	}
	return 1
}

func applyOp(vals *[]*big.Int, ops *[]byte) error {
	if len(*ops) == 0 || len(*vals) < 2 {
		return fmt.Errorf("not enough operands")
	}
	op := (*ops)[len(*ops)-1]
	*ops = (*ops)[:len(*ops)-1]
	b := (*vals)[len(*vals)-1]
	a := (*vals)[len(*vals)-2]
	*vals = (*vals)[:len(*vals)-2]

	res := new(big.Int).Set(a)
	switch op {
	case '+':
		res.Add(res, b)
	case '-':
		res.Sub(res, b)
	case '*':
		res.Mul(res, b)
	default:
		return fmt.Errorf("unknown operator %q", op)
	}
	*vals = append(*vals, res)
	return nil
}

func reverseString(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
