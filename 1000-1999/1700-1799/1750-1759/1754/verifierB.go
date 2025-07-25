package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func minDiff(arr []int) int {
	m := 1<<31 - 1
	for i := 0; i < len(arr)-1; i++ {
		d := arr[i+1] - arr[i]
		if d < 0 {
			d = -d
		}
		if d < m {
			m = d
		}
	}
	return m
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	tokens := strings.Fields(strings.TrimSpace(out.String()))
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	seen := make([]bool, n+1)
	arr := make([]int, n)
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("invalid integer %q", t)
		}
		if v < 1 || v > n || seen[v] {
			return fmt.Errorf("invalid permutation")
		}
		seen[v] = true
		arr[i] = v
	}
	md := minDiff(arr)
	if md != n/2 {
		return fmt.Errorf("expected min diff %d got %d", n/2, md)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if err := runCase(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
