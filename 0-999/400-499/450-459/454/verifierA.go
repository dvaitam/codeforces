package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n int) string {
	var sb strings.Builder
	mid := n / 2
	for i := 0; i < n; i++ {
		d := i - mid
		if d < 0 {
			d = -d
		}
		stars := d
		ds := n - 2*d
		for j := 0; j < stars; j++ {
			sb.WriteByte('*')
		}
		for j := 0; j < ds; j++ {
			sb.WriteByte('D')
		}
		for j := 0; j < stars; j++ {
			sb.WriteByte('*')
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
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
		var n int
		fmt.Sscan(line, &n)
		input := fmt.Sprintf("%d\n", n)
		exp := expected(n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("Test %d wrong answer\nExpected:\n%s\nGot:\n%s\n", idx, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
