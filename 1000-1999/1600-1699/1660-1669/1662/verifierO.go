package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type caseData struct {
	ops []string
}

type testCase struct {
	desc      string
	input     string
	caseCount int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierO.go /path/to/binary")
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
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.caseCount)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		ans, err := parseAnswers(out, tc.caseCount)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for idx := 0; idx < tc.caseCount; idx++ {
			if ans[idx] != refAns[idx] {
				fmt.Printf("Wrong answer on test %d (%s) case %d: expected %s got %s\nInput:\n%s", i+1, tc.desc, idx+1, refAns[idx], ans[idx], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1662O.bin"
	cmd := exec.Command("go", "build", "-o", path, "1662O.go")
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

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d verdicts, got %d", expected, len(tokens))
	}
	ans := make([]string, expected)
	for i, tok := range tokens {
		upper := strings.ToUpper(tok)
		if upper != "YES" && upper != "NO" {
			return nil, fmt.Errorf("invalid verdict %q", tok)
		}
		ans[i] = upper
	}
	return ans, nil
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, cases []caseData) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(cases))
		for _, cs := range cases {
			fmt.Fprintf(&sb, "%d\n", len(cs.ops))
			for _, op := range cs.ops {
				sb.WriteString(op)
				if !strings.HasSuffix(op, "\n") {
					sb.WriteByte('\n')
				}
			}
		}
		tests = append(tests, testCase{
			desc:      desc,
			input:     sb.String(),
			caseCount: len(cases),
		})
	}

	add("no-obstacles", []caseData{{ops: nil}})

	fullyBlocked := make([]caseData, 1)
	var ops []string
	for r := 1; r <= 20; r++ {
		ops = append(ops,
			fmt.Sprintf("C %d 0 180", r),
			fmt.Sprintf("C %d 180 0", r),
		)
	}
	fullyBlocked[0] = caseData{ops: ops}
	add("fully-blocked", fullyBlocked)

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 50 {
		numCases := rng.Intn(4) + 1
		cases := make([]caseData, numCases)
		for i := 0; i < numCases; i++ {
			opCnt := rng.Intn(70)
			curOps := make([]string, 0, opCnt)
			for j := 0; j < opCnt; j++ {
				if rng.Intn(2) == 0 {
					r := rng.Intn(20) + 1
					t1 := rng.Intn(360)
					t2 := rng.Intn(360)
					if t1 == t2 {
						t2 = (t2 + 1) % 360
					}
					curOps = append(curOps, fmt.Sprintf("C %d %d %d", r, t1, t2))
				} else {
					r1 := rng.Intn(20)
					r2 := rng.Intn(20-r1) + r1 + 1
					theta := rng.Intn(360)
					curOps = append(curOps, fmt.Sprintf("S %d %d %d", r1, r2, theta))
				}
			}
			cases[i] = caseData{ops: curOps}
		}
		add(fmt.Sprintf("random-%d", len(tests)), cases)
	}

	// Stress tests with many operations
	heavyCases := make([]caseData, 2)
	for idx := range heavyCases {
		curOps := make([]string, 0, 150)
		for j := 0; j < 150; j++ {
			if j%2 == 0 {
				r := (j % 20) + 1
				t1 := (j * 7) % 360
				t2 := (t1 + 173) % 360
				if t2 == t1 {
					t2 = (t2 + 1) % 360
				}
				curOps = append(curOps, fmt.Sprintf("C %d %d %d", r, t1, t2))
			} else {
				r1 := j % 20
				r2 := (r1 + 5)
				if r2 > 20 {
					r2 = 20
				}
				if r2 == r1 {
					r2 = r1 + 1
				}
				theta := (j * 13) % 360
				curOps = append(curOps, fmt.Sprintf("S %d %d %d", r1, r2, theta))
			}
		}
		heavyCases[idx] = caseData{ops: curOps}
	}
	add("heavy", heavyCases)

	return tests
}
