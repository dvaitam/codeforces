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

func expected(n, m, k int64) string {
	var ans int64
	if k == 1 {
		ans = 1
	} else if k == 2 {
		part1 := m
		if n < m {
			part1 = n
		}
		var part2 int64
		if m >= n {
			part2 = m/n - 1
		}
		if part2 < 0 {
			part2 = 0
		}
		ans = part1 + part2
	} else if k == 3 {
		part1 := m
		if n < m {
			part1 = n
		}
		var part2 int64
		if m >= n {
			part2 = m/n - 1
		}
		if part2 < 0 {
			part2 = 0
		}
		ans = m - (part1 + part2)
	} else {
		ans = 0
	}
	return fmt.Sprintf("%d", ans)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "case %d: expected 3 numbers got %d\n", idx, len(fields))
			os.Exit(1)
		}
		n, _ := strconv.ParseInt(fields[0], 10, 64)
		m, _ := strconv.ParseInt(fields[1], 10, 64)
		k, _ := strconv.ParseInt(fields[2], 10, 64)
		expectedOut := expected(n, m, k)
		input := fmt.Sprintf("1\n%d %d %d\n", n, m, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if out != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expectedOut, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
