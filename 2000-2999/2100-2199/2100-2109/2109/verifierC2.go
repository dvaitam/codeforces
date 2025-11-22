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
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
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
		// Sanity check reference output is valid.
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, gotOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2109C2.go",
		filepath.Join("2000-2999", "2100-2199", "2100-2109", "2109", "2109C2.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 2109C2.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2109C2_%d.bin", time.Now().UnixNano()))
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

func validateOutput(tc testCase, out string) error {
	tokens := strings.Fields(out)
	pos := 0

	nextToken := func() (string, bool) {
		if pos >= len(tokens) {
			return "", false
		}
		v := tokens[pos]
		pos++
		return v, true
	}

	reader := strings.NewReader(tc.input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return fmt.Errorf("bad input header: %v", err)
	}
	if t != tc.t {
		return fmt.Errorf("internal error: test case count mismatch")
	}
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var start, target int64
		if _, err := fmt.Fscan(reader, &start, &target); err != nil {
			return fmt.Errorf("failed to read pair for case %d: %v", caseIdx+1, err)
		}
		x := start
		cmdCount := 0
		for {
			cmd, ok := nextToken()
			if !ok {
				return fmt.Errorf("case %d: ran out of output before '!'", caseIdx+1)
			}
			if cmd == "!" {
				if x != target {
					return fmt.Errorf("case %d: final value %d does not equal target %d", caseIdx+1, x, target)
				}
				break
			}
			cmdCount++
			if cmdCount > 4 {
				return fmt.Errorf("case %d: used more than 4 commands", caseIdx+1)
			}

			switch cmd {
			case "add":
				valTok, ok := nextToken()
				if !ok {
					return fmt.Errorf("case %d: missing parameter for add", caseIdx+1)
				}
				y, err := strconv.ParseInt(valTok, 10, 64)
				if err != nil {
					return fmt.Errorf("case %d: invalid add parameter %q", caseIdx+1, valTok)
				}
				if y < -1e18 || y > 1e18 {
					return fmt.Errorf("case %d: add parameter out of range", caseIdx+1)
				}
				res := x + y
				if res >= 1 && res <= 1e18 {
					x = res
				}
			case "mul":
				valTok, ok := nextToken()
				if !ok {
					return fmt.Errorf("case %d: missing parameter for mul", caseIdx+1)
				}
				y, err := strconv.ParseInt(valTok, 10, 64)
				if err != nil {
					return fmt.Errorf("case %d: invalid mul parameter %q", caseIdx+1, valTok)
				}
				if y < 1 || y > 1e18 {
					return fmt.Errorf("case %d: mul parameter out of range", caseIdx+1)
				}
				if x != 0 && y > 1e18/x {
					// would overflow bounds
					break
				}
				res := x * y
				if res >= 1 && res <= 1e18 {
					x = res
				}
			case "div":
				valTok, ok := nextToken()
				if !ok {
					return fmt.Errorf("case %d: missing parameter for div", caseIdx+1)
				}
				y, err := strconv.ParseInt(valTok, 10, 64)
				if err != nil {
					return fmt.Errorf("case %d: invalid div parameter %q", caseIdx+1, valTok)
				}
				if y < 1 || y > 1e18 {
					return fmt.Errorf("case %d: div parameter out of range", caseIdx+1)
				}
				if y != 0 && x%y == 0 {
					x = x / y
				}
			case "digit":
				x = digitSum(x)
			default:
				return fmt.Errorf("case %d: unknown command %q", caseIdx+1, cmd)
			}
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra output detected after processing all cases")
	}
	return nil
}

func digitSum(x int64) int64 {
	var s int64
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 60)
	tests = append(tests, sampleTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 45; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTests() []testCase {
	var sb strings.Builder
	sb.WriteString("4\n")
	// start, target pairs reflecting various scenarios
	sb.WriteString("9 100\n")    // need addition
	sb.WriteString("1234 5\n")   // digit + div works
	sb.WriteString("7 7\n")      // no commands needed
	sb.WriteString("100 1000\n") // add big positive
	return []testCase{{input: sb.String(), t: 4}}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(6) + 1 // 1..6 cases
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		start := randRange(rng, 1, 1_000_000_000)
		target := randRange(rng, 1, 1_000_000_000)
		sb.WriteString(fmt.Sprintf("%d %d\n", start, target))
	}
	return testCase{input: sb.String(), t: t}
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}
