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

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(49) + 2 // 2..50
	used := make(map[int]bool)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(2*n) + 1
			if !used[v] {
				used[v] = true
				arr[i] = v
				break
			}
		}
	}
	var in strings.Builder
	in.WriteString("1\n")
	in.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')

	// compute expected answer
	ans := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i]*arr[j] == (i+1)+(j+1) {
				ans++
			}
		}
	}
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", ans))
	return testCase{input: in.String(), expected: out.String()}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// simple fixed case with a known pair
	cases := []testCase{{input: "1\n2\n1 3\n", expected: "1\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
