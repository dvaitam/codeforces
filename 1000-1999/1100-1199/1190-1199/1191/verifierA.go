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

func solveCase(x int) (int, byte) {
	bestRank := -1
	bestA := 0
	var bestCat byte
	for a := 0; a <= 2; a++ {
		r := (x + a) % 4
		rank := 0
		var cat byte
		switch r {
		case 1:
			rank = 4
			cat = 'A'
		case 3:
			rank = 3
			cat = 'B'
		case 2:
			rank = 2
			cat = 'C'
		case 0:
			rank = 1
			cat = 'D'
		}
		if rank > bestRank {
			bestRank = rank
			bestA = a
			bestCat = cat
		}
	}
	return bestA, bestCat
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
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
		x, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid number on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		expectA, expectCat := solveCase(x)
		input := fmt.Sprintf("%d\n", x)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		expect := fmt.Sprintf("%d %c", expectA, expectCat)
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", idx, expect, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
