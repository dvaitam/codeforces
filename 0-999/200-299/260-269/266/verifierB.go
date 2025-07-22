package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, t int, s string) string {
	b := []byte(s)
	for step := 0; step < t; step++ {
		for i := 0; i+1 < n; {
			if b[i] == 'B' && b[i+1] == 'G' {
				b[i], b[i+1] = b[i+1], b[i]
				i += 2
			} else {
				i++
			}
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	var tCases int
	fmt.Sscan(scanner.Text(), &tCases)
	for i := 0; i < tCases; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough test cases")
			os.Exit(1)
		}
		line := scanner.Text()
		var n, tt int
		var s string
		fmt.Sscan(line, &n, &tt, &s)
		input := fmt.Sprintf("%d %d\n%s\n", n, tt, s)
		want := expected(n, tt, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tCases)
}
