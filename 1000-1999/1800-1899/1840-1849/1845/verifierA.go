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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hasSolution(n, k, x int) bool {
	if x != 1 {
		return true
	}
	if k == 1 {
		return false
	}
	if n%2 == 0 {
		return true
	}
	if k == 2 {
		return false
	}
	if n == 1 {
		return false
	}
	return true
}

func checkOutput(out string, n, k, x int, expectYes bool) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(scanner.Text())
	firstUp := strings.ToUpper(first)
	if !expectYes {
		if firstUp != "NO" {
			return fmt.Errorf("expected NO got %s", first)
		}
		return nil
	}
	if firstUp != "YES" {
		return fmt.Errorf("expected YES got %s", first)
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing m line")
	}
	mLine := strings.TrimSpace(scanner.Text())
	m, err := strconv.Atoi(mLine)
	if err != nil {
		return fmt.Errorf("invalid m: %v", err)
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing numbers line")
	}
	parts := strings.Fields(scanner.Text())
	if len(parts) != m {
		return fmt.Errorf("expected %d numbers got %d", m, len(parts))
	}
	sum := 0
	for _, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("bad number %q", p)
		}
		if v < 1 || v > k || v == x {
			return fmt.Errorf("invalid value %d", v)
		}
		sum += v
	}
	if sum != n {
		return fmt.Errorf("numbers sum to %d expected %d", sum, n)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		x, _ := strconv.Atoi(parts[2])
		input := fmt.Sprintf("1\n%d %d %d\n", n, k, x)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		expectYes := hasSolution(n, k, x)
		if err := checkOutput(out, n, k, x, expectYes); err != nil {
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
