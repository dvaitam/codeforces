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

func expectedCount(members []string) int {
	base := members[0]
	cnt := 0
	for _, s := range members {
		if s == base {
			cnt++
		}
	}
	return cnt
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
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
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad testcase on line %d\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil || n <= 0 {
			fmt.Fprintf(os.Stderr, "bad n on line %d\n", idx)
			os.Exit(1)
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil || k <= 0 {
			fmt.Fprintf(os.Stderr, "bad k on line %d\n", idx)
			os.Exit(1)
		}
		if len(parts) != 2+n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d strings got %d\n", idx, n, len(parts)-2)
			os.Exit(1)
		}
		members := make([]string, n)
		for i := 0; i < n; i++ {
			members[i] = parts[2+i]
		}
		expect := expectedCount(members)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i, s := range members {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(s)
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != fmt.Sprint(expect) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
