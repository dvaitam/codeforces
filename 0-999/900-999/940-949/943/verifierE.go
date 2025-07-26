package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	nums []int
}

func solve(nums []int) string {
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)
	var b strings.Builder
	for i, v := range sorted {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

func (tc testCase) input() string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(len(tc.nums)))
	b.WriteByte('\n')
	for i, v := range tc.nums {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(v))
	}
	b.WriteByte('\n')
	return b.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(100)
	}
	return testCase{nums: nums}
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.nums)
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{randomCase(rng)}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
