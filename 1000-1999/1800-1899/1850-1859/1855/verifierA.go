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

func expectedMoves(p []int) int {
	fixed := 0
	for i, v := range p {
		if v == i+1 {
			fixed++
		}
	}
	moves := fixed / 2
	if fixed%2 == 1 {
		moves++
	}
	return moves
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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
		fmt.Fprintln(os.Stderr, "failed to open testcasesA.txt:", err)
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
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid n on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(fields)-1 != n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d: invalid number %q\n", idx, fields[i+1])
				os.Exit(1)
			}
			p[i] = val
		}
		expected := expectedMoves(p)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", p[i]))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		want := fmt.Sprintf("%d", expected)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
