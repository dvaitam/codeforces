package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(arr []int) int {
	ones := []int{}
	others := []int{}
	for _, v := range arr {
		if v == 1 {
			ones = append(ones, v)
		} else {
			others = append(others, v)
		}
	}
	sort.Ints(others)
	cnt := map[int]int{}
	for _, v := range others {
		cnt[v]++
	}
	uniq := []int{}
	for v := range cnt {
		uniq = append(uniq, v)
	}
	sort.Ints(uniq)
	seq := make([]int, 0, len(arr))
	seq = append(seq, ones...)
	for _, v := range uniq {
		seq = append(seq, v)
	}
	extras := []int{}
	for _, v := range uniq {
		for i := 1; i < cnt[v]; i++ {
			extras = append(extras, v)
		}
	}
	sort.Ints(extras)
	seq = append(seq, extras...)

	prev := 0
	total := 0
	for _, v := range seq {
		next := prev + v - (prev % v)
		if next <= prev {
			next += v
		}
		total += next
		prev = next
	}
	return total
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(6))
	tests := []testCase{}
	tests = append(tests, testCase{input: "1\n", expected: "1"})
	tests = append(tests, testCase{input: "3\n1 2 3\n", expected: fmt.Sprint(solveCase([]int{1, 2, 3}))})
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(6) + 1
		}
		res := solveCase(arr)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{input: sb.String(), expected: fmt.Sprint(res)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
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
