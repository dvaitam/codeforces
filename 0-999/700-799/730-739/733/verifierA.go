package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(s string) string {
	vowels := map[byte]bool{'A': true, 'E': true, 'I': true, 'O': true, 'U': true, 'Y': true}
	last := 0
	ans := 0
	for i := 0; i < len(s); i++ {
		if vowels[s[i]] {
			gap := i + 1 - last
			if gap > ans {
				ans = gap
			}
			last = i + 1
		}
	}
	if len(s)+1-last > ans {
		ans = len(s) + 1 - last
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
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
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("unexpected EOF")
			os.Exit(1)
		}
		s := scan.Text()
		input := fmt.Sprintf("%s\n", s)
		exp := solveCase(s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
