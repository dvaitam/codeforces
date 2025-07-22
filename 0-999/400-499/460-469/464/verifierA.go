package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveA(n, p int, s string) string {
	a := []byte(s)
	for i := n - 1; i >= 0; i-- {
		for c := a[i] + 1; c < byte('a'+p); c++ {
			if i >= 1 && a[i-1] == c {
				continue
			}
			if i >= 2 && a[i-2] == c {
				continue
			}
			a[i] = c
			ok := true
			for j := i + 1; j < n; j++ {
				found := false
				for d := byte('a'); d < byte('a'+p); d++ {
					if j >= 1 && a[j-1] == d {
						continue
					}
					if j >= 2 && a[j-2] == d {
						continue
					}
					a[j] = d
					found = true
					break
				}
				if !found {
					ok = false
					break
				}
			}
			if ok {
				return string(a)
			}
		}
	}
	return "NO"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var n, p int
		fmt.Sscan(line, &n, &p)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing string on case %d\n", idx+1)
			os.Exit(1)
		}
		s := strings.TrimSpace(scanner.Text())
		idx++
		expected := solveA(n, p, s)
		input := fmt.Sprintf("%d %d\n%s\n", n, p, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
