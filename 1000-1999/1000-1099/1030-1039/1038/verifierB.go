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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func verify(n int, out string) error {
	out = strings.TrimSpace(out)
	if out == "" {
		return fmt.Errorf("empty output")
	}
	lines := strings.Split(out, "\n")
	first := strings.TrimSpace(lines[0])
	if strings.EqualFold(first, "No") {
		if n > 2 {
			return fmt.Errorf("expected Yes for n=%d", n)
		}
		if len(lines) > 1 {
			return fmt.Errorf("unexpected extra output after No")
		}
		return nil
	}
	if !strings.EqualFold(first, "Yes") {
		return fmt.Errorf("first line must be Yes or No, got %q", first)
	}
	if n <= 2 {
		return fmt.Errorf("expected No for n=%d but got Yes", n)
	}
	if len(lines) < 3 {
		return fmt.Errorf("expected three lines for Yes output")
	}
	aFields := strings.Fields(lines[1])
	bFields := strings.Fields(lines[2])
	if len(aFields) == 0 || len(bFields) == 0 {
		return fmt.Errorf("missing subset lines")
	}
	k, err := strconv.Atoi(aFields[0])
	if err != nil {
		return fmt.Errorf("invalid subset size: %v", err)
	}
	m, err := strconv.Atoi(bFields[0])
	if err != nil {
		return fmt.Errorf("invalid subset size: %v", err)
	}
	if len(aFields)-1 != k || len(bFields)-1 != m {
		return fmt.Errorf("subset size mismatch")
	}
	if k+m != n {
		return fmt.Errorf("subsets do not partition 1..n")
	}
	used := make([]bool, n+1)
	sumA, sumB := 0, 0
	for _, tok := range aFields[1:] {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("invalid number in first subset: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("number %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("number %d repeated", v)
		}
		used[v] = true
		sumA += v
	}
	for _, tok := range bFields[1:] {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("invalid number in second subset: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("number %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("number %d repeated", v)
		}
		used[v] = true
		sumB += v
	}
	for i := 1; i <= n; i++ {
		if !used[i] {
			return fmt.Errorf("number %d missing from partition", i)
		}
	}
	if gcd(sumA, sumB) <= 1 {
		return fmt.Errorf("sums %d and %d are coprime", sumA, sumB)
	}
	return nil
}

func runCase(exe, input string, n int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verify(n, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d\n", n)
		if err := runCase(exe, input, n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
