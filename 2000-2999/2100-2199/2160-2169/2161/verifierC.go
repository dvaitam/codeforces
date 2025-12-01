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

const refSource = "./2161C.go"

type testCase struct {
	n int
	x int64
	a []int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	solutions, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	if err := verifyCandidate(candOut, tests, solutions); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\ncandidate output:\n%s", err, candOut)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_2161C.bin"
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
	return stdout.String(), nil
}

type solution struct {
	bonus int64
	order []int
}

func parseOutput(out string, tests []testCase) ([]solution, error) {
	lines := strings.Fields(out)
	ptr := 0
	res := make([]solution, len(tests))
	for i, tc := range tests {
		if ptr >= len(lines) {
			return nil, fmt.Errorf("test %d: missing bonus line", i+1)
		}
		bonus, err := strconv.ParseInt(lines[ptr], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid bonus %q", i+1, lines[ptr])
		}
		ptr++
		if ptr+tc.n > len(lines) {
			return nil, fmt.Errorf("test %d: insufficient elements for order", i+1)
		}
		order := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.ParseInt(lines[ptr+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid element %q", i+1, lines[ptr+j])
			}
			order[j] = int(val)
		}
		ptr += tc.n
		res[i] = solution{bonus: bonus, order: order}
	}
	if ptr != len(lines) {
		return nil, fmt.Errorf("extra output detected")
	}
	return res, nil
}

func verifyCandidate(out string, tests []testCase, sol []solution) error {
	cand, err := parseOutput(out, tests)
	if err != nil {
		return err
	}
	for i, tc := range tests {
		ref := sol[i]
		c := cand[i]
		if c.bonus != ref.bonus {
			return fmt.Errorf("test %d: bonus mismatch (expected %d, got %d)", i+1, ref.bonus, c.bonus)
		}
		if err := checkOrder(tc, c.order, ref.bonus); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	return nil
}

func checkOrder(tc testCase, order []int, expected int64) error {
	if len(order) != tc.n {
		return fmt.Errorf("order length mismatch: expected %d got %d", tc.n, len(order))
	}
	// verify that multiset matches
	cnt := make(map[int64]int)
	for _, v := range tc.a {
		cnt[v]++
	}
	for _, v := range order {
		val := int64(v)
		cnt[val]--
		if cnt[val] < 0 {
			return fmt.Errorf("order contains invalid element %d", v)
		}
	}
	for _, c := range cnt {
		if c != 0 {
			return fmt.Errorf("order multiset mismatch")
		}
	}

	var sum int64
	var bonus int64
	curLevel := int64(0)
	for _, val := range order {
		sum += int64(val)
		newLevel := sum / tc.x
		if newLevel > curLevel {
			curLevel = newLevel
			bonus += int64(val)
		}
	}
	if bonus != expected {
		return fmt.Errorf("order yields bonus %d (expected %d)", bonus, expected)
	}
	return nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(n int, x int64, arr []int64) {
		tests = append(tests, testCase{n: n, x: x, a: append([]int64(nil), arr...)})
	}

	add(1, 1, []int64{1})
	add(2, 1, []int64{1, 1})
	add(3, 2, []int64{1, 1, 1})
	add(5, 10, []int64{2, 2, 2, 2, 5})
	add(11, 23, []int64{5, 5, 2, 2, 1, 2, 1, 2, 5, 3, 10})
	add(1, 2, []int64{1})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(8) + 1
		x := int64(rng.Intn(10) + 1)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(int(x)) + 1)
		}
		add(n, x, arr)
	}

	// add some larger cases
	tests = append(tests, testCase{
		n: 13,
		x: 100,
		a: []int64{44, 32, 1, 16, 100, 50, 42, 80, 73, 11, 29, 9, 25},
	})

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
