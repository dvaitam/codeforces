package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(arr []int) int64 {
	positions := []int{}
	for i, v := range arr {
		if v == 1 {
			positions = append(positions, i)
		}
	}
	if len(positions) == 0 {
		return 0
	}
	ans := int64(1)
	for i := 1; i < len(positions); i++ {
		diff := positions[i] - positions[i-1]
		ans *= int64(diff)
	}
	return ans
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		hasOne := false
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 1 {
				arr[i] = 1
				hasOne = true
			}
		}
		if !hasOne {
			arr[rng.Intn(n)] = 1
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expected := fmt.Sprintf("%d", solve(arr))
		tests = append(tests, test{sb.String(), expected})
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
