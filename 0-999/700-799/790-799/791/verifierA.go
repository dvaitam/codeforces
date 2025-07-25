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

func expected(a, b int) string {
	years := 0
	for a <= b {
		a *= 3
		b *= 2
		years++
	}
	return fmt.Sprintf("%d\n", years)
}

func runCase(bin string, a, b int) error {
	input := fmt.Sprintf("%d %d\n", a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected(a, b))
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
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
		a, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scanner.Text())
		if err := runCase(bin, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
