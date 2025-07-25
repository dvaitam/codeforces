package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func possible(s string) bool {
	n := len(s)
	if n < 26 {
		return false
	}
	arr := []rune(s)
	for i := 0; i+26 <= n; i++ {
		var freq [26]int
		q := 0
		ok := true
		for j := i; j < i+26; j++ {
			ch := arr[j]
			if ch == '?' {
				q++
			} else {
				idx := ch - 'A'
				if idx < 0 || idx >= 26 {
					ok = false
					break
				}
				if freq[idx] > 0 {
					ok = false
					break
				}
				freq[idx]++
			}
		}
		if !ok {
			continue
		}
		used := 0
		for k := 0; k < 26; k++ {
			if freq[k] > 0 {
				used++
			}
		}
		if used+q == 26 {
			return true
		}
	}
	return false
}

func validOutput(orig, cand string) bool {
	if cand == "-1" {
		return false
	}
	if len(orig) != len(cand) {
		return false
	}
	for i := 0; i < len(orig); i++ {
		if orig[i] != '?' && orig[i] != cand[i] {
			return false
		}
		if cand[i] < 'A' || cand[i] > 'Z' {
			return false
		}
	}
	n := len(cand)
	if n < 26 {
		return false
	}
	for i := 0; i+26 <= n; i++ {
		var freq [26]int
		ok := true
		for j := i; j < i+26; j++ {
			idx := cand[j] - 'A'
			if idx < 0 || idx >= 26 {
				ok = false
				break
			}
			if freq[idx] > 0 {
				ok = false
				break
			}
			freq[idx]++
		}
		if !ok {
			continue
		}
		valid := true
		for k := 0; k < 26; k++ {
			if freq[k] != 1 {
				valid = false
				break
			}
		}
		if valid {
			return true
		}
	}
	return false
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesB.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		possibleFlag := possible(s)
		out, err := run(bin, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if !possibleFlag {
			if out != "-1" {
				fmt.Printf("test %d failed: expected -1 got %s\n", idx, out)
				os.Exit(1)
			}
			continue
		}
		if !validOutput(s, out) {
			fmt.Printf("test %d failed: invalid output %q\n", idx, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
