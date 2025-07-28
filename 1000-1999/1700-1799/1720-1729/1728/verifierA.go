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
	input    string
	expected string
}

func solveCase(cnt []int) int {
	pos := 1
	maxVal := cnt[0]
	for i := 1; i < len(cnt); i++ {
		if cnt[i] > maxVal || (cnt[i] == maxVal && i+1 > pos) {
			maxVal = cnt[i]
			pos = i + 1
		}
	}
	return pos
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	tests := []testCase{}
	// some fixed cases
	tests = append(tests, testCase{input: "1\n1\n5\n", expected: "1"})
	tests = append(tests, testCase{input: "1\n3\n1 1 1\n", expected: "3"})
	tests = append(tests, testCase{input: "1\n3\n1 2 3\n", expected: "3"})
	tests = append(tests, testCase{input: "1\n4\n1 5 3 5\n", expected: "4"})
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		cnt := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			cnt[i] = rng.Intn(100) + 1
			sum += cnt[i]
		}
		if sum%2 == 0 {
			cnt[0]++
			sum++
		}
		pos := solveCase(cnt)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d\n", n))
		for i, v := range cnt {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		tests = append(tests, testCase{input: sb.String(), expected: fmt.Sprint(pos)})
	}
	return tests
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	_ = rng // quiet lint
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
