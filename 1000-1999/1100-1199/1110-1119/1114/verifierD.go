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

const testcasesRaw = `4 5 5 2 3
10 4 5 1 5 1 4 3 5 2 2
8 5 5 4 4 2 2 2 5
7 1 1 2 5 1 3 1
5 4 5 4 4 4
10 4 2 3 1 1 2 4 2 3 4
5 4 5 4 5 3
9 5 4 5 2 3 1 3 5 2
6 5 5 5 1 2 5
5 3 1 1 4 4
2 3 1
7 2 1 3 4 4 1 1
10 5 1 4 5 3 5 3 5 2 1
5 1 1 1 5 5
1 2
7 3 5 3 2 1 3 3
6 2 4 4 4 5 4
10 5 1 5 5 3 4 2 3 4 3
9 3 5 3 1 4 5 3 1 4
10 5 2 1 3 4 3 3 5 3 4
1 5
1 1
6 3 4 3 5 5 3
3 3 2 3
6 5 3 3 4 1 1
10 2 3 5 2 3 2 3 2 4 1
2 5 3
6 2 4 2 1 3 2
10 4 3 2 1 1 5 2 3 5 2
5 3 1 5 3 5
3 4 3 5
5 4 3 4 3 4
10 4 1 4 2 2 1 4 5 5 4
9 2 1 4 5 3 5 3 2 1
10 3 1 2 1 1 5 2 4 5 1
1 4
2 2 5
5 2 1 5 5 4
1 5
2 3 2
5 5 4 1 3 2
4 1 5 1 2
4 3 2 1 4
10 4 1 3 2 3 5 5 5 4 1
8 3 1 1 2 1 1 1 1
8 1 1 5 5 4 3 2 3
2 3 4
7 5 3 3 3 2 3 4
2 2 5
1 4
2 5 2
1 3
8 5 5 4 1 5 4 1 3
8 3 4 4 4 1 2 2 5
5 5 1 4 2 4
3 1 3 3
9 3 1 4 1 5 4 1 3 5
9 1 5 1 4 2 2 4 1 5
2 5 1
7 2 1 3 1 1 1 4
5 5 3 1 1 5
9 5 2 1 5 1 5 1 5 3
10 2 1 2 2 2 4 5 4 3 3
10 4 3 5 4 1 4 5 2 4 2
7 5 5 5 4 2 4 2
3 1 4 4
9 4 5 2 2 3 2 2 5 5
6 2 5 3 4 5 5
10 3 2 3 1 3 4 4 2 2 5
6 2 3 4 2 4 4
10 2 4 5 5 1 4 1 4 1 4
4 2 1 2 3
4 2 3 2 2
10 1 3 2 1 3 2 4 1 1 1
2 3 3
1 3
8 5 3 1 1 3 3 4 4
8 1 2 5 4 4 2 5 3
2 3 1
7 1 4 5 3 1 5 3
6 4 3 3 1 3 5
9 5 1 4 5 3 1 3 5 2
3 2 3 4
2 1 5
3 3 5 4
9 3 2 4 4 3 2 1 1 2
9 5 5 4 3 1 3 3 4 1
3 1 4 5
5 2 5 3 3 4
8 5 1 3 4 1 2 3 5
2 1 5
2 2 2
10 4 4 2 5 5 2 4 2 5 5
3 5 2 2
5 3 3 1 4 4
7 3 5 5 3 4 5 3
8 1 5 2 1 1 2 4 2
9 4 2 2 5 2 1 5 4 1
10 3 2 2 4 1 5 1 1 3 5
4 5 1 4 5`

// referenceSolve mirrors 1114D.go.
func referenceSolve(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	v := make([]int, 0, len(arr))
	prev := arr[0]
	v = append(v, prev)
	for i := 1; i < len(arr); i++ {
		if arr[i] != prev {
			v = append(v, arr[i])
			prev = arr[i]
		}
	}
	freq := make(map[int]int)
	maxF := 0
	for _, val := range v {
		freq[val]++
		if freq[val] > maxF {
			maxF = freq[val]
		}
	}
	return len(v) - maxF
}

func runCase(bin string, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return fmt.Errorf("invalid test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid n: %v", err)
	}
	if len(fields) != 1+n {
		return fmt.Errorf("expected %d numbers got %d", 1+n, len(fields))
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i], err = strconv.Atoi(fields[1+i])
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
	}
	expect := referenceSolve(arr)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(line + "\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
