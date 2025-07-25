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

func calc(x, length int64) int64 {
	if length <= x-1 {
		return (x - 1 + x - length) * length / 2
	}
	return (x-1)*x/2 + (length - (x - 1))
}

func feasible(n, m, k, x int64) bool {
	left := calc(x, k-1)
	right := calc(x, n-k)
	return x+left+right <= m
}

func expected(n, m, k int64) string {
	l, r := int64(1), m
	var ans int64
	for l <= r {
		mid := (l + r) / 2
		if feasible(n, m, k, mid) {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(exe, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
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
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		m64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		k64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		input := fmt.Sprintf("%d %d %d\n", n64, m64, k64)
		exp := expected(n64, m64, k64) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
