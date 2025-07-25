package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expectedA(x1, x2, x3 int) int {
	xs := []int{x1, x2, x3}
	sort.Ints(xs)
	m := xs[1]
	return abs(xs[0]-m) + abs(xs[1]-m) + abs(xs[2]-m)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
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
		x1, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x2, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x3, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d %d\n", x1, x2, x3)
		exp := fmt.Sprintf("%d\n", expectedA(x1, x2, x3))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
