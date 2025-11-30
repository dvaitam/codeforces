package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `5 6 0 4 8 7 6
4 7 5 9 3 8
3 4 2 1 9
4 8 9 2 4 1
2 10 5 7
2 5 6 5
3 8 7 7 8
4 0 8 0 1 6
2 9 7 5
3 5 1 3 9
3 3 2 8 7
2 1 5 8
5 1 4 8 4 1 8
4 8 3 9 8 9
4 7 1 9 6 5
3 4 2 3 2
2 9 10 4
5 1 1 10 2 2 0
2 8 10 6
4 8 3 3 10 9
5 9 4 7 7 10 10
4 1 5 9 1 7
4 3 3 0 4 1
3 5 2 5 6
2 1 2 3
2 9 10 8
2 0 1 10
3 9 9 1 6
2 5 1 0
2 3 2 1
5 3 0 10 0 8 6
2 4 1 3
2 10 4 5
5 2 0 8 7 0 9
2 6 3 4
4 7 9 2 10 3
2 10 2 2
4 8 4 1 9 7
3 0 7 10 6
4 10 5 6 10 4
3 8 0 7 1
4 0 8 4 2 3
5 5 9 4 10 5 9
3 4 6 6 10
2 0 9 3
4 2 3 3 10 7
5 10 9 6 0 6 9
5 10 0 2 7 1 4
3 7 8 7 8
2 0 7 5
4 7 0 6 3 8
2 2 0 6
5 5 0 3 0 0 10
2 3 1 9
3 4 4 2 1
5 6 10 1 0 4 7
2 4 2 10
4 1 2 4 0 0
2 3 10 4
4 5 9 0 9 10
5 10 7 10 6 5 8
3 3 6 9 4
2 2 2 4
4 5 5 1 5 9
2 0 4 2
3 9 4 5 6
3 4 1 7 3
2 4 2 8
2 4 6 5
4 6 1 1 8 7
5 5 5 1 7 1 7
5 0 4 5 10 2 2
5 10 1 1 1 3 3
2 6 0 1
5 8 8 4 7 7 9
3 6 1 5 3
4 9 2 6 3 5
2 1 0 8
5 10 3 1 7 6 4
3 10 0 3 9
3 1 3 7 6
4 8 2 1 9 7
3 9 6 10 10
5 8 7 10 5 7 7
3 8 9 3 0
4 5 5 0 8 2
4 9 2 6 9 4
5 1 1 8 0 1 3
3 0 4 0 7
4 2 2 10 7 5
5 8 8 0 9 1 10
2 6 3 4
5 7 6 9 9 3 0
2 2 4 8
4 5 1 7 4 4
5 6 6 0 2 10 2
3 4 5 0 0
5 6 2 7 9 1 10
3 5 6 0 9
5 6 7 0 1 7 2`

// solveCase replicates the 1516A solver for a single test.
func solveCase(n, k int, arr []int) string {
	for i := 0; i < n-1 && k > 0; i++ {
		if arr[i] > 0 {
			move := arr[i]
			if move > k {
				move = k
			}
			arr[i] -= move
			arr[n-1] += move
			k -= move
		}
	}
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseLine(line string) (int, int, []int, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("bad n: %w", err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("bad k: %w", err)
	}
	if len(fields) != 2+n {
		return 0, 0, nil, fmt.Errorf("expected %d array entries, got %d", n, len(fields)-2)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("bad a[%d]: %w", i, err)
		}
		arr[i] = v
	}
	return n, k, arr, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcasesA, "\n")
	idx := 0
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		idx++
		n, k, arr, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\ninput:\n%s\n", idx, err, line)
			os.Exit(1)
		}

		input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, strings.Join(strings.Fields(line)[2:], " "))
		exp := solveCase(n, k, append([]int(nil), arr...))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
