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

func expected(n int, s string) string {
	cnt := 0
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			cnt++
		}
	}
	return fmt.Sprintf("%d", cnt)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	tline := strings.TrimSpace(scanner.Text())
	var t int
	fmt.Sscan(tline, &t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough test cases")
			os.Exit(1)
		}
		line := scanner.Text()
		var n int
		var s string
		fmt.Sscan(line, &n, &s)
		input := fmt.Sprintf("%d\n%s\n", n, s)
		want := expected(n, s)
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
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", scanner.Err())
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
