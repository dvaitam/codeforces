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
)

const refSource = "2000-2999/2100-2199/2120-2129/2124/2124D.go"

type testCase struct {
	n   int
	k   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[len(args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
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
	exp, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if exp[i] != got[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "wrong answer on case %d (n=%d k=%d): expected %v got %v\n", i+1, tc.n, tc.k, yn(exp[i]), yn(got[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func yn(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2124D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func buildTests() []testCase {
	var tcs []testCase
	add := func(n, k int, arr []int) {
		tcs = append(tcs, testCase{n: n, k: k, arr: arr})
	}

	// Deterministic cases including samples.
	add(5, 3, []int{5, 4, 3, 4, 5})
	add(4, 1, []int{1, 1, 2, 1})
	add(6, 6, []int{2, 3, 4, 5, 3, 2})
	add(5, 4, []int{5, 2, 4, 3, 1})
	add(8, 5, []int{4, 7, 1, 2, 3, 1, 3, 4})
	add(5, 4, []int{1, 2, 1, 2, 2})
	add(3, 2, []int{1, 2, 2})
	add(4, 4, []int{2, 1, 2, 2})

	rng := rand.New(rand.NewSource(2124))
	totalN := 0
	for len(tcs) < 40 && totalN < 180000 {
		n := rng.Intn(4000) + 1
		if totalN+n > 180000 {
			n = 180000 - totalN
		}
		if n == 0 {
			break
		}
		k := rng.Intn(n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		add(n, k, arr)
		totalN += n
	}

	// Large stress case close to limit but still safe for runtime.
	largeN := 200000
	largeArr := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		largeArr[i] = rng.Intn(largeN) + 1
	}
	add(largeN, largeN/2, largeArr)

	return tcs
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.Grow(32 + len(tcs)*32)
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
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

func parseAnswers(out string, t int) ([]bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]bool, t)
	for i, tok := range tokens {
		l := strings.ToLower(tok)
		if len(l) == 0 {
			return nil, fmt.Errorf("empty token at position %d", i+1)
		}
		switch l[0] {
		case 'y':
			res[i] = true
		case 'n':
			res[i] = false
		default:
			return nil, fmt.Errorf("invalid answer %q at position %d", tok, i+1)
		}
	}
	return res, nil
}
