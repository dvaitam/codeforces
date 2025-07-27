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

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func generateTests() ([]int, string) {
	const t = 100
	r := rand.New(rand.NewSource(1))
	nums := make([]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(99) + 2
		nums[i] = n
		fmt.Fprintf(&sb, "%d\n", n)
	}
	return nums, sb.String()
}

func verify(nums []int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, n := range nums {
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				return fmt.Errorf("case %d: not enough numbers", idx+1)
			}
			var x int
			fmt.Sscan(scanner.Text(), &x)
			arr[i] = x
		}
		seen := make([]bool, n+1)
		for i, v := range arr {
			if v < 1 || v > n {
				return fmt.Errorf("case %d: value out of range", idx+1)
			}
			if seen[v] {
				return fmt.Errorf("case %d: duplicate value %d", idx+1, v)
			}
			if v == i+1 {
				return fmt.Errorf("case %d: fixed point at %d", idx+1, i+1)
			}
			seen[v] = true
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	nums, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(nums, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem A")
}
