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

func solveA(a, b int64) string {
	if a < b {
		return "-1"
	}
	k := (a + b) / (2 * b)
	ans := float64(a+b) / (2 * float64(k))
	return fmt.Sprintf("%.12f", ans)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid test case format on line %d\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		expected := solveA(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		// allow small floating point error
		if expected == "-1" {
			if got != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\n", idx, got)
				os.Exit(1)
			}
			continue
		}
		valExp, _ := strconv.ParseFloat(expected, 64)
		valGot, err := strconv.ParseFloat(got, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: non-float output %s\n", idx, got)
			os.Exit(1)
		}
		if diff := valExp - valGot; diff > 1e-6 || diff < -1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
