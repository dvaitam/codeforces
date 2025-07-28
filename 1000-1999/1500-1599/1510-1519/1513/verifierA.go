package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct{ n, k int }

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest(rng *rand.Rand) (string, []testCase) {
	t := rng.Intn(20) + 1
	cases := make([]testCase, t)
	var sb strings.Builder
	fmt.Fprintln(&sb, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		k := rng.Intn(n + 1)
		cases[i] = testCase{n, k}
		fmt.Fprintf(&sb, "%d %d\n", n, k)
	}
	return sb.String(), cases
}

func checkCase(n, k int, line string) error {
	line = strings.TrimSpace(line)
	if n <= 2*k {
		if line != "-1" {
			return fmt.Errorf("expected -1")
		}
		return nil
	}
	fields := strings.Fields(line)
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	used := make([]bool, n+1)
	arr := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("not integer: %v", err)
		}
		if v < 1 || v > n || used[v] {
			return fmt.Errorf("not a permutation")
		}
		used[v] = true
		arr[i] = v
	}
	peaks := 0
	for i := 1; i < n-1; i++ {
		if arr[i] > arr[i-1] && arr[i] > arr[i+1] {
			peaks++
		}
	}
	if peaks != k {
		return fmt.Errorf("expected %d peaks got %d", k, peaks)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, cases := genTest(rng)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != len(cases) {
			fmt.Printf("test %d: expected %d lines got %d\ninput:\n%s\n", i+1, len(cases), len(lines), input)
			os.Exit(1)
		}
		for j, cs := range cases {
			if err := checkCase(cs.n, cs.k, lines[j]); err != nil {
				fmt.Printf("test %d case %d failed: %v\ninput:\n%s\n", i+1, j+1, err, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
