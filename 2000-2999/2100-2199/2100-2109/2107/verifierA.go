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

const refSource = "./2107A.go"

type testCase struct {
	n    int
	arr  []int
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refFeasible, err := parseRef(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateCandidate(candOut, tests, refFeasible); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2107A.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, arr []int) {
		cpy := append([]int(nil), arr...)
		tests = append(tests, testCase{n: len(cpy), arr: cpy, name: name})
	}

	add("two_diff", []int{1, 2})
	add("two_same", []int{5, 5})
	add("all_same", []int{7, 7, 7, 7})
	add("simple_mix", []int{1, 2, 3})
	add("large_vals", []int{10000, 9999, 10000, 5000})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	const maxTotalN = 40000
	for len(tests) < 200 && totalN < maxTotalN {
		n := rng.Intn(100) + 2 // 2..101
		if totalN+n > maxTotalN {
			break
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			switch rng.Intn(4) {
			case 0:
				arr[i] = rng.Intn(10) + 1
			case 1:
				arr[i] = rng.Intn(10000) + 1
			case 2:
				arr[i] = arr[0]
			default:
				arr[i] = rng.Intn(50) + 1
			}
		}
		add(fmt.Sprintf("random_%d", len(tests)), arr)
		totalN += n
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseRef(out string, tests []testCase) ([]bool, error) {
	tokens := strings.Fields(out)
	pos := 0
	res := make([]bool, len(tests))
	for i := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("reference output ended early at case %d", i+1)
		}
		ans := strings.ToUpper(tokens[pos])
		pos++
		if ans == "YES" {
			res[i] = true
		} else if ans == "NO" {
			res[i] = false
		} else {
			return nil, fmt.Errorf("case %d: expected YES/NO got %q", i+1, tokens[pos-1])
		}
		if res[i] {
			need := tests[i].n
			if pos+need > len(tokens) {
				return nil, fmt.Errorf("case %d: expected %d group values, got %d", i+1, need, len(tokens)-pos)
			}
			pos += need // ignore actual assignment
		}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("reference output has %d extra tokens", len(tokens)-pos)
	}
	return res, nil
}

func validateCandidate(out string, tests []testCase, feasible []bool) error {
	tokens := strings.Fields(out)
	pos := 0
	for i, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("candidate output ended early at case %d", i+1)
		}
		ans := strings.ToUpper(tokens[pos])
		pos++
		if ans == "NO" {
			if feasible[i] {
				return fmt.Errorf("case %d (%s): reference says YES but candidate says NO", i+1, tc.name)
			}
			continue
		}
		if ans != "YES" {
			return fmt.Errorf("case %d (%s): expected YES/NO got %q", i+1, tc.name, tokens[pos-1])
		}
		if !feasible[i] {
			return fmt.Errorf("case %d (%s): reference says NO but candidate says YES", i+1, tc.name)
		}
		if pos+tc.n > len(tokens) {
			return fmt.Errorf("case %d (%s): expected %d group values, got %d", i+1, tc.name, tc.n, len(tokens)-pos)
		}
		assign := make([]int, tc.n)
		seen1, seen2 := false, false
		for j := 0; j < tc.n; j++ {
			v, err := strconv.Atoi(tokens[pos+j])
			if err != nil || (v != 1 && v != 2) {
				return fmt.Errorf("case %d (%s): invalid assignment %q at position %d", i+1, tc.name, tokens[pos+j], j+1)
			}
			assign[j] = v
			if v == 1 {
				seen1 = true
			} else {
				seen2 = true
			}
		}
		pos += tc.n
		if !seen1 || !seen2 {
			return fmt.Errorf("case %d (%s): both groups must be non-empty", i+1, tc.name)
		}
		if !checkGCD(tc.arr, assign) {
			return fmt.Errorf("case %d (%s): gcds are not different", i+1, tc.name)
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("candidate output has %d extra tokens", len(tokens)-pos)
	}
	return nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func checkGCD(arr []int, assign []int) bool {
	var g1, g2 int
	first1, first2 := true, true
	for i, v := range arr {
		if assign[i] == 1 {
			if first1 {
				g1 = v
				first1 = false
			} else {
				g1 = gcd(g1, v)
			}
		} else {
			if first2 {
				g2 = v
				first2 = false
			} else {
				g2 = gcd(g2, v)
			}
		}
	}
	return g1 != g2
}
