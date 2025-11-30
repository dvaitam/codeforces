package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one line: n followed by n values).
const embeddedTestcases = `2 2 1
3 3 2 1
7 3 5 1 2 6 7 4
1 1
7 6 2 3 1 5 7 4
3 2 1 3
5 2 4 5 3 1
1 1
6 2 6 3 4 1 5
4 2 4 3 1
7 3 1 6 2 5 7 4
6 5 6 1 3 4 2
7 7 4 3 2 1 6 5
1 1
6 2 1 6 4 5 3
3 1 3 2
7 3 6 2 1 7 4 5
6 1 3 5 2 6 4
1 1
4 3 2 4 1
3 3 1 2
4 4 3 2 1
6 4 3 1 5 6 2
5 1 3 4 2 5
7 2 1 7 6 4 3 5
4 4 1 3 2
4 4 3 2 1
5 1 4 3 2 5
4 1 4 2 3
1 1
5 4 1 2 3 5
7 3 4 1 5 7 6 2
1 1
7 3 4 2 5 7 1 6
7 5 4 1 6 2 3 7
3 3 1 2
2 2 1
3 2 1 3
6 1 2 5 4 6 3
4 2 3 1 4
3 1 3 2
4 3 1 4 2
3 2 1 3
5 2 3 5 1 4
4 4 3 1 2
4 2 1 3 4
6 2 3 1 4 5 6
4 1 2 4 3
1 1
6 6 4 1 5 2 3
1 1
7 7 5 4 2 6 3 1
2 2 1
5 5 4 3 2 1
2 2 1
4 3 1 4 2
2 1 2
5 3 5 1 4 2
3 1 2 3
1 1
3 1 2 3
3 2 3 1
2 1 2
7 7 6 1 4 3 2 5
6 6 2 5 3 4 1
7 6 3 4 7 1 2 5
2 2 1
5 1 4 5 3 2
3 1 3 2
1 1
3 3 2 1
2 2 1
3 3 2 1
1 1
4 3 4 1 2
1 1
5 2 3 1 4 5
5 5 4 2 3 1
5 4 3 2 1 5
1 1
7 6 2 4 5 7 3 1
1 1
7 6 4 3 5 2 1 7
2 2 1
4 4 1 3 2
2 2 1
4 1 2 3 4
5 5 1 2 4 3
2 1 2
1 1
1 1
1 1
7 1 4 2 7 5 6 3
3 3 1 2
1 1
3 1 2 3
1 1
3 3 2 1
6 4 1 5 2 6 3
3 2 3 1`

// Embedded solver logic from 351B.go.
func solve(vals []int) int64 {
	var inversions int64
	n := len(vals)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if vals[j] < vals[i] {
				inversions++
			}
		}
	}
	return (inversions/2)*4 + inversions%2
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil || n < 0 || len(parts) != 1+n {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx+1)
			os.Exit(1)
		}
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			vals[i], _ = strconv.Atoi(parts[1+i])
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range vals {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := strconv.FormatInt(solve(vals), 10)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
