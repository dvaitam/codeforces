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

func expectedB(n int, s string) (int, int) {
	inside := false
	current := 0
	maxOut := 0
	inWords := 0
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			current++
		} else {
			if current > 0 {
				if inside {
					inWords++
				} else if current > maxOut {
					maxOut = current
				}
				current = 0
			}
			if ch == '(' {
				inside = true
			} else if ch == ')' {
				inside = false
			}
		}
	}
	if current > 0 {
		if inside {
			inWords++
		} else if current > maxOut {
			maxOut = current
		}
	}
	return maxOut, inWords
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		expOut, expIn := expectedB(n, s)
		input := fmt.Sprintf("%d\n%s\n", n, s)
		exp := fmt.Sprintf("%d %d\n", expOut, expIn)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
