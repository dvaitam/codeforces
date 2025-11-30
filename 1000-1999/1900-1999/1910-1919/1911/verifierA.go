package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `9 14 14 14 14 2 14 14 14 14
10 13 13 13 13 13 13 13 10 13 13
8 19 19 7 19 19 19 19 19
7 5 5 5 5 4 5 5
7 18 20 18 18 18 18 18
7 4 4 4 4 4 4 3
8 16 18 16 16 16 16 16 16
8 14 14 14 11 14 14 14 14
10 15 15 15 15 17 15 15 15 15 15
3 1 18 18
9 1 1 1 1 1 1 1 20 1
8 8 11 8 8 8 8 8 8
6 19 8 19 19 19 19
5 15 18 18 18 18
4 11 11 11 17
4 10 10 18 10
4 18 11 18 18
7 15 15 15 15 3 15 15
9 11 11 11 19 11 11 11 11 11
7 6 6 6 6 6 6 7
5 2 2 20 2 2
10 3 3 5 3 3 3 3 3 3 3
3 3 3 18
9 17 17 17 17 17 17 17 17 9
6 7 7 7 19 7 7
7 15 15 15 15 15 16 15
8 3 11 3 3 3 3 3 3
10 19 19 19 11 19 19 19 19 19 19
6 9 1 1 1 1 1
6 12 12 6 12 12 12
9 2 2 4 2 2 2 2 2 2
6 2 2 2 2 2 19
4 1 4 1 1
4 13 13 3 13
4 20 2 2 2
6 6 6 6 4 6 6
6 2 2 2 2 1 2
9 20 20 20 20 4 20 20 20 20
4 8 8 3 8
8 6 14 14 14 14 14 14 14
10 2 20 2 2 2 2 2 2 2 2
9 7 7 7 7 7 9 7 7 7
10 19 19 19 6 19 19 19 19 19 19
3 6 6 11
7 4 4 4 20 4 4 4
5 1 1 1 16 1
7 12 12 12 12 12 12 13
7 5 5 5 5 5 18 5
3 15 3 15
3 9 18 18
6 16 16 16 16 12 16
7 12 12 12 12 12 19 12
5 10 10 10 13 10
4 1 20 1 1
8 6 6 6 8 6 6 6 6
10 13 13 13 13 13 13 19 13 13 13
3 13 19 13
3 15 6 6
7 6 6 6 6 15 6 6
10 20 18 18 18 18 18 18 18 18 18
3 16 11 16
10 2 2 2 14 2 2 2 2 2 2
4 5 5 5 1
9 11 11 11 1 11 11 11 11 11
3 1 1 17
4 7 4 7 7
7 6 9 9 9 9 9 9
10 3 13 13 13 13 13 13 13 13 13
7 15 15 15 15 15 15 4
7 5 5 5 5 5 5 17
8 4 4 4 4 5 4 4 4
3 2 2 7
7 18 18 11 18 18 18 18
3 20 20 16
10 14 14 14 14 14 14 14 14 12 14
5 7 7 7 7 13
7 1 5 1 1 1 1 1
7 11 11 11 11 11 12 11
4 20 11 11 11
3 6 9 9
7 12 12 12 12 13 12 12
5 10 10 10 4 10
6 2 10 2 2 2 2
4 10 10 13 10
7 4 14 14 14 14 14 14
10 16 16 16 16 16 11 16 16 16 16
4 16 16 16 4
9 2 2 2 2 2 10 2 2 2
5 6 6 6 19 6
4 3 7 3 3
3 1 13 13
9 18 18 18 18 17 18 18 18 18
10 16 16 16 19 16 16 16 16 16 16
9 3 3 3 12 3 3 3 3 3
7 19 19 19 6 19 19 19
6 4 12 12 12 12 12
3 17 17 15
6 4 4 4 16 4 4
7 7 7 7 7 7 7 2
6 5 20 20 20 20 20`

func runBinary(bin string, input string) (string, error) {
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

func solve(arr []int) int {
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	for i, v := range arr {
		if freq[v] == 1 {
			return i + 1
		}
	}
	return -1
}

func parseCases(raw string) ([][]int, error) {
	lines := strings.Split(raw, "\n")
	var res [][]int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n: %w", err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("length mismatch n=%d got %d", n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid value: %w", err)
			}
			arr[i] = v
		}
		res = append(res, arr)
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, arr := range cases {
		want := solve(arr)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		outStr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("Test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(outStr))
		if err != nil {
			fmt.Printf("Test %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
