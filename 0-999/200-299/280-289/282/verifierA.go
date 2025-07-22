package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(ops []string) int {
	x := 0
	for _, op := range ops {
		if strings.Contains(op, "+") {
			x++
		} else {
			x--
		}
	}
	return x
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	n := 0
	fmt.Sscan(fields[0], &n)
	if len(fields)-1 != n {
		return fmt.Errorf("test %d: expected %d ops got %d", idx, n, len(fields)-1)
	}
	ops := fields[1:]
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", n)
	for _, op := range ops {
		input.WriteString(op)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var ans int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &ans); err != nil {
		return fmt.Errorf("could not parse output: %v", err)
	}
	exp := expected(ops)
	if ans != exp {
		return fmt.Errorf("expected %d got %d", exp, ans)
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
		fmt.Fprintf(os.Stderr, "cannot open testcasesA.txt: %v\n", err)
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
