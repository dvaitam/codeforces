package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	arr     []int
	input   string
	maxVals []int
	minVals []int
}

func solveCase(arr []int) ([]int, []int) {
	n := len(arr)
	positives := 0
	negatives := 0
	for _, v := range arr {
		if v > 0 {
			positives++
		} else {
			negatives++
		}
	}
	maxVals := make([]int, n)
	for i := 1; i <= n; i++ {
		if i <= positives {
			maxVals[i-1] = i
		} else {
			maxVals[i-1] = 2*positives - i
		}
	}
	minVals := make([]int, n)
	paired := negatives * 2
	for i := 1; i <= n; i++ {
		if i <= paired {
			if i%2 == 1 {
				minVals[i-1] = 1
			} else {
				minVals[i-1] = 0
			}
		} else {
			minVals[i-1] = i - paired
		}
	}
	return maxVals, minVals
}

func buildCase(arr []int) testCase {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	maxVals, minVals := solveCase(arr)
	return testCase{arr: arr, input: sb.String(), maxVals: maxVals, minVals: minVals}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	negatives := rng.Intn(n/2 + 1)
	positives := n - negatives
	arr := make([]int, n)
	idx := 0
	for i := 0; i < positives; i++ {
		arr[idx] = rng.Intn(n) + 1
		idx++
	}
	for i := 0; i < negatives; i++ {
		arr[idx] = -(rng.Intn(n) + 1)
		idx++
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return buildCase(arr)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkOutput(out string, tc testCase) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		return fmt.Errorf("expected 2 lines got %d", len(lines))
	}
	fields1 := strings.Fields(lines[0])
	fields2 := strings.Fields(lines[1])
	if len(fields1) != len(tc.maxVals) || len(fields2) != len(tc.minVals) {
		return fmt.Errorf("expected %d numbers per line", len(tc.maxVals))
	}
	for i, f := range fields1 {
		var val int
		if _, err := fmt.Sscan(f, &val); err != nil || val != tc.maxVals[i] {
			return fmt.Errorf("expected %v got %v", tc.maxVals, fields1)
		}
	}
	for i, f := range fields2 {
		var val int
		if _, err := fmt.Sscan(f, &val); err != nil || val != tc.minVals[i] {
			return fmt.Errorf("expected %v got %v", tc.minVals, fields2)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// simple deterministic cases
	cases = append(cases, buildCase([]int{1}))
	cases = append(cases, buildCase([]int{1, -1}))
	cases = append(cases, buildCase([]int{1, 2, -1, -2}))
	cases = append(cases, buildCase([]int{1, 2, 3}))
	cases = append(cases, buildCase([]int{1, -1, 2}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d mismatch: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
