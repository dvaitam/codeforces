package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(s string) string {
	stack := make([]rune, 0, len(s))
	cnt := 0
	for _, ch := range s {
		n := len(stack)
		if n > 0 && stack[n-1] == ch {
			stack = stack[:n-1]
			cnt++
		} else {
			stack = append(stack, ch)
		}
	}
	if cnt%2 == 1 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		input := s + "\n"
		exp := expected(s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%s\noutput:\n%s\n", idx, err, input, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(strings.ToUpper(out))
		if out != exp {
			fmt.Printf("test %d failed for input %q\nexpected: %s\noutput: %s\n", idx, s, exp, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
