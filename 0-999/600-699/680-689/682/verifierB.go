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
	arr []int
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func solve(arr []int) int {
	sort.Ints(arr)
	mex := 1
	for _, v := range arr {
		if v >= mex {
			mex++
		}
	}
	return mex
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
	for i, v := range tc.arr {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < len(tc.arr) {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := solve(append([]int(nil), tc.arr...))
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(200) + 1
	}
	return testCase{arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{arr: []int{1}})
	cases = append(cases, testCase{arr: []int{1, 2, 3, 3, 4}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, tc.arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
