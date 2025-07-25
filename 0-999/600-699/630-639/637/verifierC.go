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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hamming(a, b string) int {
	d := 0
	for i := 0; i < 6; i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func expected(codes []string) int {
	if len(codes) == 1 {
		return 6
	}
	minDist := 7
	for i := 0; i < len(codes); i++ {
		for j := i + 1; j < len(codes); j++ {
			d := hamming(codes[i], codes[j])
			if d < minDist {
				minDist = d
				if minDist == 1 {
					break
				}
			}
		}
		if minDist == 1 {
			break
		}
	}
	if minDist < 1 {
		minDist = 1
	}
	return (minDist - 1) / 2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty testcases file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test count")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", caseIdx)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || n != len(fields)-1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		codes := fields[1:]
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(codes, " "))
		want := fmt.Sprintf("%d", expected(codes))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseIdx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
