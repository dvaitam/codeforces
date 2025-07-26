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

func solve(L, v, l, r int64) int64 {
	total := L / v
	inRange := r/v - ((l - 1) / v)
	return total - inRange
}

func runCase(bin string, L, v, l, r int64) error {
	input := fmt.Sprintf("1\n%d %d %d %d\n", L, v, l, r)
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
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", solve(L, v, l, r))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("invalid test case on line %d\n", idx+1)
			os.Exit(1)
		}
		L, _ := strconv.ParseInt(parts[0], 10, 64)
		v, _ := strconv.ParseInt(parts[1], 10, 64)
		l, _ := strconv.ParseInt(parts[2], 10, 64)
		r, _ := strconv.ParseInt(parts[3], 10, 64)
		if err := runCase(bin, L, v, l, r); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
