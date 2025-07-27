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

func expected(rs, bs string) string {
	red, blue := 0, 0
	for i := 0; i < len(rs); i++ {
		if rs[i] > bs[i] {
			red++
		} else if rs[i] < bs[i] {
			blue++
		}
	}
	switch {
	case red > blue:
		return "RED"
	case red < blue:
		return "BLUE"
	default:
		return "EQUAL"
	}
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
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
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("bad test case %d\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("bad n in test %d\n", idx)
			os.Exit(1)
		}
		rs, bs := parts[1], parts[2]
		if len(rs) != n || len(bs) != n {
			fmt.Printf("test %d length mismatch\n", idx)
			os.Exit(1)
		}
		expectedOut := expected(rs, bs)
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, rs, bs)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expectedOut, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
