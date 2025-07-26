package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(n int, m int64, heights []int64) int {
	sort.Slice(heights, func(i, j int) bool { return heights[i] > heights[j] })
	ans := 0
	i := 0
	for heights[0] != heights[n-1] {
		t := i
		var tot int64
		for t+1 < n {
			diff := heights[t] - heights[t+1]
			cost := diff * int64(t+1)
			if tot+cost > m {
				break
			}
			tot += cost
			t++
		}
		if t == i {
			i = t
		}
		ans++
		i = t
		dec := (m - tot) / int64(i+1)
		heights[i] -= dec
	}
	return ans
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(44))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		m := rng.Int63n(1000) + int64(n)
		heights := make([]int64, n)
		for i := 0; i < n; i++ {
			heights[i] = rng.Int63n(200) + 1
		}
		hcopy := make([]int64, n)
		copy(hcopy, heights)
		ans := solveCase(n, m, hcopy)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range heights {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{input: sb.String(), expected: fmt.Sprint(ans)})
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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
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
