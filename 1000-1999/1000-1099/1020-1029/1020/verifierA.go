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

func expected(n, h, a, b, n1, h1, n2, h2 int) int {
	if n1 == n2 {
		if h2 > h1 {
			return h2 - h1
		}
		return h1 - h2
	}
	bsv := h1
	if h1 > b {
		bsv = b
	} else if h1 < a {
		bsv = a
	}
	diff := h2 - bsv
	if diff < 0 {
		diff = -diff
	}
	diff2 := bsv - h1
	if diff2 < 0 {
		diff2 = -diff2
	}
	diff3 := n2 - n1
	if diff3 < 0 {
		diff3 = -diff3
	}
	return diff + diff2 + diff3
}

func runCase(exe string, n, h, a, b, n1, h1, n2, h2 int) error {
	input := fmt.Sprintf("%d %d %d %d 1\n%d %d %d %d\n", n, h, a, b, n1, h1, n2, h2)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strconv.Itoa(expected(n, h, a, b, n1, h1, n2, h2))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
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
		nums := make([]int, 8)
		for j := 0; j < 8; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			nums[j], _ = strconv.Atoi(scan.Text())
		}
		if err := runCase(exe, nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], nums[6], nums[7]); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
