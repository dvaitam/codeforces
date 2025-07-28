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

func runCandidate(bin, input string) (string, error) {
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

func expected(n int, s1, s2 string) string {
	counts := make(map[[2]byte]int)
	b1 := []byte(s1)
	b2 := []byte(s2)
	for i := 0; i < n; i++ {
		a := b1[i]
		b := b2[n-1-i]
		if a > b {
			a, b = b, a
		}
		counts[[2]byte{a, b}]++
	}
	oddSame := 0
	ok := true
	for pair, c := range counts {
		if pair[0] != pair[1] {
			if c%2 == 1 {
				ok = false
				break
			}
		} else {
			if c%2 == 1 {
				oddSame++
			}
		}
	}
	if ok {
		if n%2 == 0 {
			if oddSame > 0 {
				ok = false
			}
		} else {
			if oddSame != 1 {
				ok = false
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesD.txt: %v\n", err)
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
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		s1 := parts[1]
		s2 := parts[2]
		expect := expected(n, s1, s2)
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, s1, s2)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
