package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	desc  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		got, err := parseAnswer(out)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Wrong answer on test %d (%s): expected %d got %d\nInput:\n%s", i+1, tc.desc, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1431F.bin"
	cmd := exec.Command("go", "build", "-o", path, "1431F.go")
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

func parseAnswer(out string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var ans int64
	if _, err := fmt.Fscan(reader, &ans); err != nil {
		return 0, fmt.Errorf("failed to read integer answer: %v", err)
	}
	if extra := strings.TrimSpace(readRemaining(reader)); extra != "" {
		return 0, fmt.Errorf("unexpected extra output: %q", extra)
	}
	return ans, nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, n, k, x int, arr []int) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, k, x)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{desc: desc, input: sb.String()})
	}

	add("single-element", 1, 1, 1, []int{5})
	add("remove-all", 3, 3, 2, []int{4, 4, 4})
	add("no-removal-necessary", 5, 1, 5, []int{1, 2, 3, 4, 5})
	add("x-equals-1", 4, 2, 1, []int{10, 1, 10, 1})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 60 {
		n := rng.Intn(50) + 1
		k := rng.Intn(n) + 1
		x := rng.Intn(n) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(1000) + 1
		}
		desc := fmt.Sprintf("random-small-%d", len(tests))
		add(desc, n, k, x, arr)
	}

	// Large stress test
	n := 100000
	k := 50000
	x := 12345
	arr := make([]int, n)
	for i := range arr {
		arr[i] = (i % 1000) + 1
	}
	add("large-structured", n, k, x, arr)

	n2 := 100000
	k2 := 100000
	x2 := 1
	arr2 := make([]int, n2)
	for i := range arr2 {
		arr2[i] = rng.Intn(100000) + 1
	}
	add("large-x1", n2, k2, x2, arr2)

	return tests
}
