package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(a, b string) string {
	if len(a) != len(b) {
		return "NO"
	}
	hasA1 := strings.Contains(a, "1")
	hasB1 := strings.Contains(b, "1")
	if hasA1 == hasB1 {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return fmt.Errorf("test %d: need two strings", idx)
	}
	a := fields[0]
	b := fields[1]
	input := a + "\n" + b + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	ans := strings.TrimSpace(strings.ToUpper(out.String()))
	exp := expected(a, b)
	if ans != exp {
		return fmt.Errorf("expected %s got %s", exp, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open testcasesC.txt: %v\n", err)
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
		idx++
		if err := runCase(bin, line, idx); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
