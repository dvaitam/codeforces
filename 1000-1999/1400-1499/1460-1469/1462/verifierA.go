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

func solve(b []int) string {
	n := len(b)
	kOdd := (n + 1) / 2
	kEven := n / 2
	a := make([]int, n)
	for i := 0; i < kOdd; i++ {
		a[2*i] = b[i]
	}
	for j := 0; j < kEven; j++ {
		idx := 2*(kEven-j) - 1
		a[idx] = b[kOdd+j]
	}
	out := make([]string, n)
	for i, v := range a {
		out[i] = fmt.Sprint(v)
	}
	return strings.Join(out, " ")
}

func runCase(bin string, n int, arr []int) error {
	// build input for binary
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintln(n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}

	got := strings.TrimSpace(out.String())
	exp := solve(arr)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
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
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scanner.Text())
			arr[j] = v
		}
		if err := runCase(bin, n, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
